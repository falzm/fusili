package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/falzm/fusili/kvmap"
	"github.com/falzm/fusili/logger"
	"github.com/falzm/fusili/output"
	"github.com/falzm/fusili/portscan"
)

type hostPort struct {
	host string
	port int
}

const (
	defaultLogLevel        string = "info"
	defaultScanConcurrency int    = 10
	defaultScanPorts       string = "1:1024"
	defaultScanPortTimeout int    = 1
)

var (
	version   string
	buildDate string

	flagScanConcurrency int
	flagScanPorts       string
	flagScanPortTimeout int
	flagConfigPath      string
	flagHelp            bool
	flagLogLevel        string
	flagVersion         bool

	logLevel int
)

func init() {
	var err error

	flag.StringVar(&flagConfigPath, "c", "", "path to configuration file")
	flag.BoolVar(&flagHelp, "h", false, "display this help and exit")
	flag.StringVar(&flagLogLevel, "l", defaultLogLevel, "logging level (error, warning, notice, info, debug)")
	flag.BoolVar(&flagVersion, "v", false, "display version and exit")
	flag.IntVar(&flagScanConcurrency, "sC", defaultScanConcurrency, "scan concurrency")
	flag.StringVar(&flagScanPorts, "sP", defaultScanPorts, "scan ports range (START:END)")
	flag.IntVar(&flagScanPortTimeout, "sT", defaultScanPortTimeout, "scan port timeout (in seconds)")
	flag.Usage = func() { printUsage(os.Stderr) }

	flag.Parse()

	if logLevel, err = logger.GetLevelByName(flagLogLevel); err != nil {
		dieOnError("invalid log level %q\n", flagLogLevel)
	}

	logger.Init(os.Stdout, true)
	logger.SetLevel(logLevel)
}

func main() {
	var (
		workers   []chan struct{}
		scanners  map[string]*portscan.PortScanner
		hostPorts map[string][]int
		outputs   map[string]output.Output
		scan      chan hostPort
		sink      chan hostPort
		wg        sync.WaitGroup
		err       error
	)

	if flagHelp {
		printUsage(os.Stdout)
		os.Exit(0)
	} else if flagVersion {
		printVersion(version, buildDate)
		os.Exit(0)
	}

	startPort, endPort, err := parseScanPorts(flagScanPorts)
	if err != nil {
		logger.Error("core", "invalid scan ports range")
		os.Exit(1)
	}

	// Load configuration from file
	config, err := loadConfig(flagConfigPath)
	if err != nil {
		logger.Error("core", "unable to load configuration: %s", err)
		os.Exit(1)
	}

	// Setup outputs
	outputs = make(map[string]output.Output)
	for outputName, outputSettings := range config.Output {
		settings, ok := outputSettings.(map[string]interface{})
		if !ok {
			logger.Error("core", "unable to initalize output %q: invalid settings structure", outputName)
			os.Exit(1)
		}

		outputType, err := kvmap.GetString(settings, "type", true)
		if err != nil {
			logger.Error("core", "unable to initalize output %q: %s", outputName, err)
			os.Exit(1)
		}

		if _, present := output.Outputs[outputType]; !present {
			logger.Error("core", "unable to initalize output %q: unsupported output type %q", outputName, outputType)
			os.Exit(1)
		}

		if outputs[outputName], err = output.Outputs[outputType](outputName, settings); err != nil {
			logger.Error("core", "unable to initalize output %q: %s", outputName, err)
			os.Exit(1)
		}
	}

	if len(config.Hosts) == 0 {
		logger.Error("core", "no hosts to scan")
		os.Exit(1)
	}

	if len(outputs) == 0 {
		logger.Error("core", "no outputs configured")
		os.Exit(1)
	}

	workers = make([]chan struct{}, flagScanConcurrency)
	hostPorts = make(map[string][]int)
	scanners = make(map[string]*portscan.PortScanner)
	scan = make(chan hostPort)
	sink = make(chan hostPort)
	wg = sync.WaitGroup{}

	// Start workers
	for w := range workers {
		workers[w] = make(chan struct{})

		go func(scan, sink chan hostPort, shutdown chan struct{}) {
			for {
				select {
				case hp := <-scan:
					if scanners[hp.host].IsOpen(strconv.Itoa(hp.port)) {
						if !scanners[hp.host].IsPortExpected(hp.port) {
							sink <- hostPort{hp.host, hp.port}
						}
					}

				case <-shutdown:
					wg.Done()
					return
				}
			}
		}(scan, sink, workers[w])

		wg.Add(1)
	}

	// Start scan results sink
	go func() {
		var hp hostPort

		for hp = range sink {
			if _, present := hostPorts[hp.host]; !present {
				hostPorts[hp.host] = make([]int, 0)
			}

			hostPorts[hp.host] = append(hostPorts[hp.host], hp.port)
		}
	}()

	logger.Debug("core", "starting scan")

	timeStart := time.Now()

	// Scan hosts
	for host, expectedPorts := range config.Hosts {
		if scanners[host], err = portscan.NewPortScanner(host, expectedPorts, flagScanPortTimeout); err != nil {
			logger.Error("core", "unable to resolve address for host %s: %s", host, err)
			continue
		}

		for p := startPort; p <= endPort; p++ {
			scan <- hostPort{host, p}
		}
	}

	// Shutdown workers
	for w := range workers {
		workers[w] <- struct{}{}
	}

	wg.Wait()
	close(sink)

	timeEnd := time.Now()

	// Report scan results
	for outputName, output := range outputs {
		if err := output.Report(hostPorts); err != nil {
			logger.Error("report", "%s: %s", outputName, err)
		}
	}

	logger.Info("core", "scanned %d hosts in %.1f seconds", len(config.Hosts), timeEnd.Sub(timeStart).Seconds())
}

func printUsage(output io.Writer) {
	fmt.Fprintf(output, "Usage: %s [OPTIONS]", path.Base(os.Args[0]))
	fmt.Fprint(output, "\n\nOptions:\n")

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(output, "   -%s  %s\n", f.Name, f.Usage)
	})

	os.Exit(2)
}

func printVersion(version, buildDate string) {
	fmt.Printf("%s version %s, built on %s\nGo version: %s (%s)\n",
		path.Base(os.Args[0]),
		version,
		buildDate,
		runtime.Version(),
		runtime.Compiler,
	)
}

func dieOnError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("error: %s", format), a)
	os.Exit(1)
}

func parseScanPorts(ports string) (int, int, error) {
	var (
		start int
		end   int
	)

	if _, err := fmt.Sscanf(ports, "%d:%d", &start, &end); err != nil {
		return 0, 0, err
	}

	return start, end, nil
}

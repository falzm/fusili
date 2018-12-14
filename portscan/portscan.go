package portscan

import (
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"github.com/falzm/fusili/logger"
)

type PortScanner struct {
	host          *net.IPAddr
	expectedPorts []int
	timeout       time.Duration
}

// NewPortScanner returns a new PortScanner instance with time out set to timeout.
func NewPortScanner(host string, expectedPorts []int, timeout int) (*PortScanner, error) {
	var err error

	ps := &PortScanner{
		expectedPorts: expectedPorts,
		timeout:       time.Duration(timeout) * time.Second,
	}

	if ps.host, err = net.ResolveIPAddr("ip", host); err != nil {
		return nil, err
	}

	return ps, nil
}

// IsOpen returns true if port on the host is open, false otherwise.
func (ps *PortScanner) IsOpen(port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ps.host.String(), port), ps.timeout)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			netError := err.(*net.OpError)

			switch reflect.TypeOf(netError.Err).String() {
			case "*os.SyscallError":
				syscallError := netError.Err.(*os.SyscallError)

				if syscallError.Err == syscall.ECONNREFUSED {
					logger.Debug("scan", "%s: port %s closed", ps.host.String(), port)
					return false
				}

			case "*net.timeoutError":
				logger.Debug("scan", "%s: port %s timed out (probably filtered by firewall)", ps.host.String(), port)
				return false
			}
		}

		logger.Error("scan", "%s: port %s: %s", ps.host.String(), port, err)
		return false
	}

	defer conn.Close()

	logger.Debug("scan", "%s: port %s open", ps.host.String(), port)

	return true
}

//
func (ps *PortScanner) IsPortExpected(port int) bool {
	for _, p := range ps.expectedPorts {
		if p == port {
			return true
		}
	}

	return false
}

// ScanRange scans the host for open ports from start to end, and returns a list of open ports.
func (ps *PortScanner) ScanRange(start, end int) []int {
	open := make([]int, 0)

	for p := start; p <= end; p++ {
		port := strconv.Itoa(p)

		log.Printf("DEBUG: ScanRange: scanning %s:%s\n", ps.host.String(), port)

		if ps.IsOpen(port) {
			open = append(open, p)
		}
	}

	return open
}

// Host returns the target host of the port scanner.
func (ps *PortScanner) Host() string {
	return ps.host.String()
}

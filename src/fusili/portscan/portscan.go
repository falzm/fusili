package portscan

import (
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"fusili/logger"
)

type PortScanner struct {
	host          string
	expectedPorts []int
	timeout       time.Duration
}

// NewPortScanner returns a new PortScanner instance with time out set to timeout.
func NewPortScanner(host string, expectedPorts []int, timeout int) *PortScanner {
	return &PortScanner{
		host:          host,
		expectedPorts: expectedPorts,
		timeout:       time.Duration(timeout) * time.Second,
	}
}

// IsOpen returns true if port on the host is open, false otherwise.
func (ps *PortScanner) IsOpen(port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ps.host, port), ps.timeout)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			netError := err.(*net.OpError)

			switch reflect.TypeOf(netError.Err).String() {
			case "*os.SyscallError":
				syscallError := netError.Err.(*os.SyscallError)

				if syscallError.Err == syscall.ECONNREFUSED {
					logger.Debug("scan", "%s: port %s closed", ps.host, port)
					return false
				}

			case "*net.timeoutError":
				logger.Debug("scan", "%s: port %s timed out (probably filtered by firewall)", ps.host, port)
				return false
			}
		}

		logger.Error("scan", "%s: port %s: %s", ps.host, port, err)
		return false
	}

	defer conn.Close()

	logger.Debug("scan", "%s: port %s open", ps.host, port)

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

		log.Printf("DEBUG: ScanRange: scanning %s:%s\n", ps.host, port)

		if ps.IsOpen(port) {
			open = append(open, p)
		}
	}

	return open
}

// Host returns the target host of the port scanner.
func (ps *PortScanner) Host() string {
	return ps.host
}

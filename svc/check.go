package svc

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type HostPort struct {
	Host  string
	Port  uint64
	PType string
}

func (hp HostPort) Check(timeout time.Duration) bool {
	if !hp.CheckPortType() {
		_, _ = fmt.Fprintf(os.Stderr, "Unknown/unsupported port type, type %s.\n", hp.PType)
		return false
	}
	hostPort := net.JoinHostPort(hp.Host, strconv.FormatUint(hp.Port, 10))
	conn, err := net.DialTimeout(hp.PType, hostPort, timeout)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Connect to port failed, error %v.\n", err)
		return false
	}
	defer conn.Close()
	return true
}

func (hp HostPort) CheckPortType() bool {
	types := []string{
		"tcp",
		"tcp4",
		"tcp6",
		"udp",
		"udp4",
		"udp6",
	}
	return searchArray(hp.PType, types) != -1
}

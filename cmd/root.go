package cmd

import (
	"flag"
	"fmt"
	"github.com/yjqg6666/port-checker/svc"
	"os"
	"strconv"
	"time"
)

const (
	AppName    = "port-checker"
	AppVersion = "v0.1.2"
)
const (
	ExitCodeOk              = 0  //man sysexits
	ExitCodeUsage           = 64 //man sysexits
	ExitCodePortUnavailable = 69
	ExitCodePortOk          = 0
)

var GType string
var GHost string
var GPort uint64
var GVerbose bool
var GVersion bool
var GTimeout uint
var GInterval uint
var GRetry uint

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: "+AppName+" host port\n")
	flag.PrintDefaults()
	os.Exit(ExitCodeUsage)
}

func defineArg() {
	flag.StringVar(&GType, "t", "tcp", "type, tcp or udp")
	flag.BoolVar(&GVerbose, "v", false, "verbose")
	flag.BoolVar(&GVersion, "version", false, "version")
	flag.UintVar(&GTimeout, "timeout", 5, "timeout in seconds")
	flag.UintVar(&GInterval, "interval", 5, "interval in seconds")
	flag.UintVar(&GRetry, "retry", 1, "retry numbers")
	flag.Usage = usage
}

func parseArg() {

	flag.Parse()

	if GVersion {
		fmt.Printf("%s\n", AppVersion)
		os.Exit(ExitCodeOk)
	}

	args := flag.Args()
	if len(args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "host and port are missing.\n")
		os.Exit(ExitCodeUsage)
	}

	host := flag.Arg(0)
	GHost = host

	port, err := strconv.ParseUint(flag.Arg(1), 10, 16)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Invalid port %q: %s\n", port, err)
		os.Exit(ExitCodeUsage)
	}
	GPort = port
}

func init() {
	defineArg()
	parseArg()
}

func RootExecute() {
	if GVerbose {
		fmt.Printf("Check host: %s %s port: %d\t", GHost, GType, GPort)
		fmt.Printf("timeout: %d, interval %d, retry: %d\n", GTimeout, GInterval, GRetry)
	}

	timeout := fmt.Sprintf("%ds", GTimeout)
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Format timeout error.\n")
		os.Exit(ExitCodeUsage)
	}

	sleep := fmt.Sprintf("%ds", GInterval)
	duration, err := time.ParseDuration(sleep)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Format interval time error.\n")
		os.Exit(ExitCodeUsage)
	}

	hp := svc.HostPort{
		Host:  GHost,
		Port:  GPort,
		PType: GType,
	}

	typeOk := hp.CheckPortType()
	if !typeOk {
		_, _ = fmt.Fprintf(os.Stderr, "Unsupported checker type %s.\n", GType)
		os.Exit(ExitCodeUsage)
	}

	for i := GRetry; i > 0; i-- {
		result := hp.Check(timeoutDuration)

		if result {
			if GVerbose {
				fmt.Printf("Check host: %s port: %d type: %s ok.\n", GHost, GPort, GType)
			}
			os.Exit(ExitCodePortOk)
		}

		left := i - 1
		if left > 0 {
			_, _ = fmt.Fprintf(os.Stderr, "Check host: %s port: %d type: %s failed, retry in %s, left %d retry.\n", GHost, GPort, GType, sleep, left)
			time.Sleep(duration)
		}
	}
	os.Exit(ExitCodePortUnavailable)
}

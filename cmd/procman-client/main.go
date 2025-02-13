package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// main sends an HTTP request to the process manager server to start, stop, or
// restart a process.
// The server is expected to be running on the host and port specified by the
// -host and -port flags.
// The command to execute is provided as the first argument.
// The server's response is logged to the console.
// The program exits with a status code of 1 if the server returns an error.
//
// Usage:
//
//	procman-client [flags] {start,stop,restart}
//	flags:
//	  -host string
//	        host to connect to (default "localhost")
//	  -port int
//	        port to connect to (default 1337)
func main() {
	var host string
	var port int

	// Parse command line flags
	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.IntVar(&port, "port", 1337, "port to connect to")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"usage:\n  %s [flags] {start,stop,restart}", filepath.Base(os.Args[0]))
	}
	flag.Parse()

	// Check if a command was provided
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Check if the command is valid
	command := flag.Arg(0)
	switch command {
	case "start", "stop", "restart":
	default:
		flag.Usage()
		os.Exit(1)
	}

	// Send an HTTP request to the server with the command
	addr := fmt.Sprintf("%s:%d/%s", host, port, command)
	resp, err := http.Post(addr, "text/plain", nil)
	if err != nil {
		log.Fatalf("http error: %v", err)
	}
	defer resp.Body.Close()

	// Check if the server returned an error
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("%s error: %v", command, resp.Status)
	}

	log.Printf("%s success", command)
}

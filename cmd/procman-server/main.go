package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/paskozdilar/procman"
)

// main starts the process manager and the HTTP server and connects
// them using a channel.
func main() {
	var port int

	// Parse command line flags
	flag.IntVar(&port, "port", 1337, "port to listen on")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"usage:\n  %[1]s [flags] COMMAND [ARGS...]:\n",
			filepath.Base(os.Args[0]))
		fmt.Fprintf(flag.CommandLine.Output(), "flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Check if a command was provided
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	log.Printf("command: %s %v\n", flag.Arg(0), flag.Args()[1:])

	// Create a channel for requests
	ChReq := make(chan procman.Request)

	// Start the process manager
	pm := &procman.ProcMan{
		Cmd:   flag.Arg(0),
		Args:  flag.Args()[1:],
		ChReq: ChReq,
	}
	go pm.Run()

	// Start the HTTP server
	pms := &procman.ProcManServer{
		ChReq: ChReq,
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("listening on port", port)
	http.Serve(l, pms)
}

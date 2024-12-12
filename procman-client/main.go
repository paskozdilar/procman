package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/paskozdilar/procman/proto_gen/procman"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var host string
	var port int

	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.IntVar(&port, "port", 8080, "port to connect to")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"usage:\n  %s [flags] {start,stop,restart}", filepath.Base(os.Args[0]))
	}
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	command := flag.Arg(0)

	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := procman.NewProcessManagerClient(conn)

	switch command {
	case "start":
		_, err := client.Start(context.Background(), &procman.StartRequest{})
		if err != nil {
			log.Fatalf("failed to start: %v", err)
		}
		log.Println("started")
	case "stop":
		_, err := client.Stop(context.Background(), &procman.StopRequest{})
		if err != nil {
			log.Fatalf("failed to stop: %v", err)
		}
		log.Println("stopped")
	case "restart":
		_, err := client.Restart(context.Background(), &procman.RestartRequest{})
		if err != nil {
			log.Fatalf("failed to restart: %v", err)
		}
		log.Println("restarted")
	}
}

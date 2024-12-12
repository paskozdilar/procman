package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"procman/proto_gen/procman"

	"google.golang.org/grpc"
)

type Action byte

const (
	ActionStart Action = iota
	ActionStop
	ActionRestart
)

type Request struct {
	Action   Action
	Response chan<- error
}

func main() {
	var port int

	// Parse command line flags
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"usage:\n  %[1]s [flags] COMMAND [ARGS...]:\n",
			filepath.Base(os.Args[0]))
		fmt.Fprintf(flag.CommandLine.Output(), "flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Check if there are any arguments
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Create a channel for requests
	ChReq := make(chan Request)

	// Start the process manager
	pm := &ProcMan{
		Cmd:   os.Args[1],
		Args:  os.Args[2:],
		ChReq: ChReq,
	}
	go pm.run()

	// Start the gRPC server
	pms := &ProcManServer{
		ChReq: ChReq,
	}
	server := grpc.NewServer()
	procman.RegisterProcessManagerServer(server, pms)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("listening on port", port)
	server.Serve(l)
}

type ProcMan struct {
	Cmd   string
	Args  []string
	ChReq chan Request
}

func (pm *ProcMan) run() {
	proc := &proc{
		Cmd:  exec.Command(pm.Cmd, pm.Args...),
		Done: make(chan error),
	}
	go proc.run()

	for {
		select {
		case req := <-pm.ChReq:
			switch req.Action {
			case ActionStart:
				if proc.Cmd != nil {
					req.Response <- errors.New("process running")
					continue
				}
				proc.Cmd = exec.Command(pm.Cmd, pm.Args...)
				proc.Done = make(chan error)
				go proc.run()
			case ActionStop:
				if proc.Cmd == nil {
					req.Response <- errors.New("process not running")
					continue
				}
				proc.Cmd.Process.Kill()
			case ActionRestart:
				if proc.Cmd == nil {
					req.Response <- errors.New("process not running")
					continue
				}
				proc.Cmd.Process.Kill()
				<-proc.Done
				proc.Cmd = exec.Command(pm.Cmd, pm.Args...)
				proc.Done = make(chan error)
				go proc.run()
			}
		case err := <-proc.Done:
			if err != nil {
				log.Println("process finished with error:", err)
			} else {
				log.Println("process finished successfully")
			}
			proc.Cmd = exec.Command(pm.Cmd, pm.Args...)
			proc.Done = make(chan error)
			go proc.run()
		}
	}
}

type proc struct {
	Cmd  *exec.Cmd
	Done chan error
}

func (p *proc) run() {
	p.Cmd.Stdin = os.Stdin
	p.Cmd.Stdout = os.Stdout
	p.Cmd.Stderr = os.Stderr

	err := p.Cmd.Start()
	if err != nil {
		p.Done <- err
		return
	}

	p.Done <- p.Cmd.Wait()
}

type ProcManServer struct {
	ChReq chan<- Request
	procman.UnimplementedProcessManagerServer
}

func (pms *ProcManServer) Start(ctx context.Context, req *procman.StartRequest) (*procman.StartResponse, error) {
	return nil, nil
}

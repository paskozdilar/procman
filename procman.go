package procman

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

// Action represents an action to be performed on the process manager.
// The possible values are:
//   - ActionStart
//   - ActionStop
//   - ActionRestart
type Action byte

const (
	ActionStart Action = iota
	ActionStop
	ActionRestart
)

// Request represents a request to be sent to the process manager.
// It contains the action to be performed and a channel to send the response
// back.
type Request struct {
	Action   Action
	Response chan<- error
}

// ProcMan represents a process manager.
// It contains the command to be executed and the arguments to be passed to it.
// It also contains a channel to receive requests.
type ProcMan struct {
	Cmd   string
	Args  []string
	ChReq chan Request
}

// Run starts the process manager.
// It listens for requests and performs the actions accordingly.
// It also listens for the process to finish and restarts it when it does.
// If the process finishes with an error, it logs the error.
// If the process finishes successfully, it logs a message.
func (pm *ProcMan) Run() {
	// Create and start the initial process.
	proc := &proc{
		Cmd:  exec.Command(pm.Cmd, pm.Args...),
		Done: make(chan error),
	}
	go proc.run()

	for {
		// Wait for one of the following events:
		// 1. A request is received
		// 2. The process finishes
		// If a request is received, perform the action accordingly.
		// If the process finishes, restart it.
		select {
		case req := <-pm.ChReq:
			switch req.Action {
			case ActionStart:
				// If the process is already running, return an error.
				if proc.Cmd != nil {
					req.Response <- errors.New("process running")
					continue
				}
				// Start the process.
				proc.Cmd = exec.Command(pm.Cmd, pm.Args...)
				proc.Done = make(chan error)
				go proc.run()
				req.Response <- nil
				log.Println("process started")
			case ActionStop:
				// If the process is not running, return an error.
				if proc.Cmd == nil {
					req.Response <- errors.New("process not running")
					continue
				}
				// Stop the process.
				proc.Cmd.Process.Kill()
				<-proc.Done
				proc.Cmd = nil
				req.Response <- nil
				log.Println("process stopped")
			case ActionRestart:
				// If the process is not running, return an error.
				if proc.Cmd == nil {
					req.Response <- errors.New("process not running")
					continue
				}
				// Restart the process.
				proc.Cmd.Process.Kill()
				<-proc.Done
				proc.Cmd = exec.Command(pm.Cmd, pm.Args...)
				proc.Done = make(chan error)
				go proc.run()
				req.Response <- nil
				log.Println("process restarted")
			}
		case err := <-proc.Done:
			// If the process finishes, restart it.
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

// proc represents a process.
// It contains the command to be executed and a channel to send the status of
// the process back.
type proc struct {
	Cmd  *exec.Cmd
	Done chan error
}

// run starts the process.
// It sets the standard input, output, and error of the process to the standard
// input, output, and error of the current process.
// It starts the process and sends the status back when it finishes.
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

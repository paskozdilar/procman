package procman

import (
	"fmt"
	"net/http"
)

// ProcManServer represents an HTTP server that receives requests to start,
// stop, or restart a process manager.
// It contains a channel to send the requests to the process manager.
// It implements the http.Handler interface.
type ProcManServer struct {
	ChReq chan<- Request
}

// Statically verify that ProcManServer implements the http.Handler interface.
var _ http.Handler = (*ProcManServer)(nil)

// ServeHTTP handles HTTP requests to start, stop, or restart a process manager.
// It reads the action from the request path and sends a request to the process
// manager.
// It writes a response with the status code 200 if the request was successful.
// If the request fails, it writes a response with the status code 500 and the
// error message.
func (pms *ProcManServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		action  Action
		command string
		chResp  chan error = make(chan error)
	)

	// Check the request path and set the action accordingly
	switch r.URL.Path {
	case "/start":
		action = ActionStart
		command = "start"
	case "/stop":
		action = ActionStop
		command = "stop"
	case "/restart":
		action = ActionRestart
		command = "restart"
	default:
		http.Error(w, "invalid action", http.StatusBadRequest)
		return
	}

	// Send the request to the process manager
	pms.ChReq <- Request{
		Action:   action,
		Response: chResp,
	}

	// Wait for the response
	if err := <-chResp; err != nil {
		http.Error(w,
			fmt.Sprintf("error: %v", err.Error()),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v success\n", command)))
}

# procman

Simple HTTP single-process manager.

## Usage

Use procman to start and stop a process on a machine. The process is started
with the given command and arguments and is restarted if it exits.

### Starting procman

To start procman, run the following command, substituting:
- ${PORT} with the port on which procman should listen for HTTP requests
- ${COMMAND} with the command to run
- ${ARGS[@]} with the list of arguments to pass to the command

```shell
procman --port=${PORT} ${COMMAND} ${ARGS[@]}
```

Procman will start the command with the given arguments and listen for HTTP
requests on the specified port.

If the process exits on its own, procman will restart it.

### Using procman

Procman listens for HTTP requests on the specified port.

Available HTTP requests are:
- `/start`: starts a stopped process with the given command and arguments
- `/stop`: stops a running process
- `/restart`: restarts a running process

### Using procman as Docker entrypoint

Procman can be used as the entrypoint for a Docker container. This allows
starting and stopping the process running in the container using HTTP requests.

To use procman in a Docker container, create a Dockerfile with the following
content:

```Dockerfile
FROM golang:latest AS procman

RUN CGO_ENABLED=0 go install github.com/paskozdilar/procman/cmd/procman-server@latest


FROM debian:bullseye

COPY --from=procman /go/bin/procman-server /go/bin/procman-server

# Your build steps

CMD ["/go/bin/procman-server", "-port=1337", "your_command", "your_command_args"]
```

Then you can start and stop service by executing HTTP command no port 1337,
e.g.:

```shell
curl http://localhost:1337/stop
```

See the [Dockerfile](./example/Dockerfile) in the example directory for a
complete example.

## Drawbacks

Procman is a simple process manager and does not provide many features that
other process managers do. 

If a process spawns child processes, procman will not manage them.
That implies that if the parent process crashes, the child processes will not be
stopped.

If the process crashes, procman will not restart it.

## Alternatives

- [supervisord](http://supervisord.org/): a more feature-rich process manager
- [daemontools](https://cr.yp.to/daemontools.html): a collection of tools for managing UNIX services
- [runit](http://smarden.org/runit/): a UNIX init scheme with service supervision

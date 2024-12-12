# procman

Simple gRPC-controlled single-process manager.

## Usage

Use procman to start and stop a process on a machine. The process is started
with the given command and arguments and is restarted if it exits.
The process can be stopped and restarted using gRPC calls.

### Starting procman

To start procman, run the following command, substituting:
- ${PORT} with the port on which procman should listen for RPC calls
- ${COMMAND} with the command to run
- ${ARGS[@]} with the list of arguments to pass to the command

```shell
procman --port=${PORT} ${COMMAND} ${ARGS[@]}
```

Procman will start the command with the given arguments and listen for RPC
calls on the specified port.

If the process exits on its own, procman will restart it.

### Using procman

Procman listens for gRPC calls on the specified port.

Available gRPC calls are:
- `StartProcess`: starts a stopped process with the given command and arguments
- `StopProcess`: stops a running process
- `RestartProcess`: stops a running process and starts it again with the same
  command and arguments

Use `proto/procman.proto` to generate the gRPC client code.

## Drawbacks

Procman is a simple process manager and does not provide many features that
other process managers do. 

If a process spawns child processes, procman will not manage them.
That implies that if the parent process crashes, the child processes will not be
stopped.

If the process crashes, procman will not restart it with a delay.

## Alternatives

- [supervisord](http://supervisord.org/): a more feature-rich process manager
- [daemontools](https://cr.yp.to/daemontools.html): a collection of tools for managing UNIX services
- [runit](http://smarden.org/runit/): a UNIX init scheme with service supervision

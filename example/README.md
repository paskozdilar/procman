# procman example

To run the example, run the following commands in this directory:

```shell
Docker build -t procman-example .
Docker run -p 1337:1337 procman-example
```

Then, you can start and stop the process by sending HTTP requests to the server
via curl:

```shell
curl http://localhost:1337/stop
curl http://localhost:1337/start
curl http://localhost:1337/restart
```



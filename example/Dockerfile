FROM golang:latest AS procman

RUN CGO_ENABLED=0 go install github.com/paskozdilar/procman/cmd/procman-server@latest


FROM debian:bullseye

COPY --from=procman /go/bin/procman-server /go/bin/procman-server

WORKDIR /app

COPY main.sh /app/main.sh

CMD ["/go/bin/procman-server", "-port=1337", "bash", "/app/main.sh"]

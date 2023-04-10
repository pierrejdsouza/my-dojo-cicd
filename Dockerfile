FROM golang:1.14-alpine

WORKDIR /src/
COPY main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/dojo
EXPOSE 8080
ENTRYPOINT ["/bin/dojo"]

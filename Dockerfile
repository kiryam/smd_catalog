FROM golang

ADD . /go/src/github.com/kiryam/smd_catalog

ENTRYPOINT go run src/github.com/kiryam/smd_catalog/cmd/main.go

# Document that the service listens on port 8080.
EXPOSE 8080
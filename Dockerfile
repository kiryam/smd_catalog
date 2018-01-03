# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/kiryam/smd_catalog
ADD . /go/src/smd_catalog

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)


#RUN curl https://glide.sh/get | sh
#RUN glide install
RUN go install github.com/kiryam/smd_catalog



# Run the outyet command by default when the container starts.
ENTRYPOINT /go/src/smd_catalog/cmd/cmd

# Document that the service listens on port 8080.
EXPOSE 8080
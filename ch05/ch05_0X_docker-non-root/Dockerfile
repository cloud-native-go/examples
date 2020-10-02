# This is a multi-stage Dockerfile. The first part executes a build in a Golang
# container, and the second retrieves the binary from the build container and
# inserts it into a "scratch" image.

# Part 1: Run all tests containerized Golang environment
#
FROM golang:1.13 as test

# Copy the source files from the host
COPY . /go/src/kvs

# Set the working directory to the same place we copied the code
WORKDIR /go/src/kvs

# Download the dependency code
RUN go mod download

# Run the tests. If the tests fail, the build will fail.
RUN go test -v .


# Part 2: Compile the binary in a containerized Golang environment
#
FROM golang:1.13 as build

# Copy the source files from the host
COPY . /go/src/kvs

# Copy the dependency source from the test container so we don't have to
# re-download it.
COPY --from=test /go/pkg/mod/ /go/pkg/mod/

# Set the working directory to the same place we copied the code
WORKDIR /go/src/kvs

# Build the binary. Note the flags that we use here:
#  CGO_ENABLED=0 --> Do not use CGO; compile statically
#  GOOS=linux    --> Compile for Linux OS
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kvs

FROM ubuntu:latest as user
RUN useradd -u 10001 kv-user
RUN touch /transactions.log && chown kv-user /transactions.log

# Part 3: Build the Key-Value Store image proper.
#
# Note that we use a "scratch" image, which contains no distribution files. The
# resulting image and containers will have only one file: our service binary.
#
FROM scratch as image

# Copy the binary from the build container.
COPY --from=build /go/src/kvs/kvs .

COPY --from=user /etc/passwd /etc/passwd
COPY --from=user /transactions.log /transactions.log

USER kv-user

# Tell Docker we'll be using port 8080
EXPOSE 8080

# Tell Docker to execute this command on a "docker run"
CMD ["/kvs"]

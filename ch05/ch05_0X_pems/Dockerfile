# This is a multi-stage Dockerfile.
#
# - The first stage executes uses "go test" to run any tests
#
# - The second stage executes compiles the Go code.
#
# - The third retrieves the binary from the build container and inserts it into a "scratch" image.

# Part 1: Run all tests containerized Golang environment
#
FROM golang:1.16 as test

# Copy the source files from the host
COPY . /src

# Set the working directory to the same place we copied the code
WORKDIR /src

# Run the tests. If the tests fail, the build will fail.
RUN go test -v .


# Part 2: Compile the binary in a containerized Golang environment
#
FROM golang:1.16 as build

# Copy the source files from the host
COPY . /src

# Copy the dependency source from the test container so we don't have to
# re-download it.
COPY --from=test /go/pkg/mod/ /go/pkg/mod/

# Set the working directory to the same place we copied the code
WORKDIR /src

# Build the binary. Note the flags that we use here:
#  CGO_ENABLED=0 --> Do not use CGO; compile statically
#  GOOS=linux    --> Compile for Linux OS
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kvs


# Part 3: Build the Key-Value Store image proper.
#
# Note that we use a "scratch" image, which contains no distribution files. The
# resulting image and containers will have only one file: our service binary.
#
FROM scratch

# Copy the binary from the build container.
COPY --from=build /src/kvs .
COPY --from=build /src/*.pem .

# Tell Docker we'll be using port 8080
EXPOSE 8080

# Tell Docker to execute this command on a "docker run"
CMD ["/kvs"]

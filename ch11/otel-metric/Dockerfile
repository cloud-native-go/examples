# Part 2: Compile the binary in a containerized Golang environment
#
FROM golang:1.15 as build

# Copy the source files from the host
COPY . /go/src/example

# Set the working directory to the same place we copied the code
WORKDIR /go/src/example

# Build the binary. Note the flags that we use here:
#  CGO_ENABLED=0 --> Do not use CGO; compile statically
#  GOOS=linux    --> Compile for Linux OS
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o example


# Part 3: Build the Key-Value Store image proper.
#
# Note that we use a "scratch" image, which contains no distribution files. The
# resulting image and containers will have only one file: our service binary.
#
FROM scratch as image

# Copy the binary from the build container.
COPY --from=build /go/src/example/example .

# Tell Docker we'll be using port 2222
EXPOSE 2222

# Tell Docker to execute this command on a "docker run"
CMD ["/example"]

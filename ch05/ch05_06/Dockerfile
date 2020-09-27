# We use a "scratch" image, which contains no distribution files. The
# resulting image and containers will have only our service binary.
FROM scratch

# Copy the existing binary from the host.
COPY kvs .

# Tell Docker we'll be using port 8080
EXPOSE 8080

# Tell Docker to execute this command on a "docker run"
CMD ["/kvs"]

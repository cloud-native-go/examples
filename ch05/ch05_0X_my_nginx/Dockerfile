# The parent image. At built time, this image will be pulled and subsequent
# instructions run against it.
FROM ubuntu:18.04

# Update apt cache and install nginx without an approval prompt.
RUN apt update && apt install --yes nginx

# Tell Docker we'll be using port 80
EXPOSE 80

# Run Nginx in the foreground. This is important: without a foreground
# process the container will automatically stop.
CMD ["nginx", "-g", "daemon off;"]

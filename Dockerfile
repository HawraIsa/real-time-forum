# Building a Docker image from the official Golang image
FROM golang

# Copying everything in the current directory to the /ascii directory in the container
COPY . /forum

# Setting the working directory inside the container
WORKDIR /forum

# Building the main.go file
RUN go build -v .

# Exposing port 8080 for the application
EXPOSE 8080

# Setting the environment variable for the port
ENV PORT=8080

# Running the application
CMD ["./forum"]

# Metadata for the created image
LABEL description="Running Ascii Art Web build using Docker"
LABEL version="1.0"
LABEL maintainer="hasan, yusuf, hawra"

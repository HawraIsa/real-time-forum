#!/bin/bash

# Name and tag for the Docker image
FULL_IMAGE_NAME="forum:latest"

# Building the Docker image
echo "Building the Docker image..."
sudo docker build -t "${FULL_IMAGE_NAME}" .
BUILD_STATUS=$?

if [ $BUILD_STATUS -eq 0 ]; then
    echo "Docker image built successfully."
else
    echo "Docker image failed to build."
    exit 1
fi

# Running the Docker container
CONTAINER_NAME="forum-container"
echo "Running the '${CONTAINER_NAME}' container using port 8080..."
sudo docker container run -p 8080:8080 --detach --name "${CONTAINER_NAME}" "${FULL_IMAGE_NAME}"
RUN_STATUS=$?

if [ $RUN_STATUS -eq 0 ]; then
    echo "Container '${CONTAINER_NAME}' is running on http://localhost:8080"
else
    echo "Failed to start the container '${CONTAINER_NAME}'."
    exit 1
fi

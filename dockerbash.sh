#!/bin/bash

# Configuration
IMAGE_NAME="forum-image"
CONTAINER_NAME="forum-container"

build_image() {
    echo "Building Docker image..."
    sudo docker build -t $IMAGE_NAME .
}

run_container() {
    echo "Running container..."
    sudo docker run -d --name $CONTAINER_NAME -p 8080:8080 $IMAGE_NAME
}

stop_container() {
    echo "Stopping container..."
    sudo docker stop $CONTAINER_NAME
}

remove_container() {
    echo "Removing container..."
    sudo docker rm $CONTAINER_NAME
}

case "$1" in
    build)
        build_image
        ;;
    run)
        run_container
        ;;
    stop)
        stop_container
        ;;
    remove)
        remove_container
        ;;
    *)
        echo "Usage: $0 {build|run|stop|remove}"
        exit 1
esac

exit 0

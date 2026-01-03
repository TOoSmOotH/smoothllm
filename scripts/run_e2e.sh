#!/bin/sh
set -e

if ! docker info >/dev/null 2>&1; then
    echo "Starting Docker daemon..."
    if command -v dockerd-entrypoint.sh >/dev/null 2>&1; then
        dockerd-entrypoint.sh &
        sleep 5
        while ! docker info >/dev/null 2>&1; do
            echo "Waiting for Docker daemon..."
            sleep 1
        done
    else
        echo "Error: Docker daemon is not running and dockerd-entrypoint.sh not found."
        exit 1
    fi
fi

echo "Docker daemon is up."

# Install docker-compose if needed (docker-compose v2 is 'docker compose')
if ! docker compose version >/dev/null 2>&1; then
    echo "Installing docker-compose plugin..."
    mkdir -p ~/.docker/cli-plugins/
    wget -O ~/.docker/cli-plugins/docker-compose https://github.com/docker/compose/releases/download/v2.23.3/docker-compose-linux-x86_64
    chmod +x ~/.docker/cli-plugins/docker-compose
fi

echo "Building and running tests..."

# We need to make sure the network is cleaned up
trap "docker compose -f docker-compose.yml -f docker-compose.test.yml down -v" EXIT

# Build images
docker compose -f docker-compose.yml build

# Run the e2e tests
# We use --abort-on-container-exit so if e2e finishes (pass or fail), everything stops
docker compose -f docker-compose.yml -f docker-compose.test.yml up --abort-on-container-exit --exit-code-from e2e

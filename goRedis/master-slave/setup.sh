#!/bin/bash

# Stop and remove existing containers
echo "Stopping and removing existing Redis containers..."
docker-compose down

# Remove any dangling images to free up space
echo "Removing unused Docker images..."
docker image prune -f

# Build new images
echo "Building Docker images..."
docker-compose build --no-cache

# Start containers in detached mode
echo "Starting Redis Master and Slaves..."
docker-compose up -d

# Wait for a few seconds to ensure services are up
echo "Waiting for Redis services to start..."
sleep 5

# Show running containers
echo "Running Docker containers:"
docker ps

# Check Master replication status
echo "Checking Redis Master replication status..."
docker exec -it master-slave-master-1 redis-cli info replication

echo "Redis Master-Slave setup is running successfully!"

docker image prune -a -f
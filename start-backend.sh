#!/bin/bash
set -e

# Check Docker is installed
if ! command -v docker &> /dev/null; then
  echo "Error: Docker is not installed. Please install Docker Desktop first."
  echo "  https://www.docker.com/products/docker-desktop/"
  exit 1
fi

# Check Docker daemon is running
if ! docker info &> /dev/null; then
  echo "Error: Docker is not running. Please start Docker Desktop and try again."
  exit 1
fi

echo "Starting backend services (MongoDB + Go server)..."
docker compose up --build

echo ""
echo "Backend is running at http://localhost:8080"
echo "API example: http://localhost:8080/api/v1/fridge"
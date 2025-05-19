#!/bin/bash

set -e

echo "Starting the application with Docker Compose..."

docker compose up --build

echo "Application is running!"
echo "API is available at http://localhost:8080"
echo "MySQL database is available at localhost:3306"
echo ""
echo "To view logs:"
echo "  docker-compose logs -f app"
echo ""
echo "To stop the application:"
echo "  docker-compose down"

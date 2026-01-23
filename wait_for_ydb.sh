#!/bin/bash
set -e

CONTAINER_PATTERN="${1:-ydb}"
MAX_WAIT="${2:-60}"
COUNTER=0

# Найти контейнер по паттерну
CONTAINER_NAME=$(docker ps --format "{{.Names}}" | grep -E "${CONTAINER_PATTERN}" | head -1)

if [ -z "$CONTAINER_NAME" ]; then
    echo "No container found matching pattern: $CONTAINER_PATTERN"
    exit 1
fi

echo "Found container: $CONTAINER_NAME"
echo "Waiting for YDB to be healthy..."

while [ $COUNTER -lt $MAX_WAIT ]; do
    STATUS=$(docker inspect --format="{{.State.Health.Status}}" "$CONTAINER_NAME" 2>/dev/null || echo "not-found")
    
    if [ "$STATUS" = "healthy" ]; then
        echo "YDB is healthy!"
        exit 0
    fi
    
    echo "YDB status: $STATUS ($COUNTER/$MAX_WAIT sec)"
    sleep 1
    COUNTER=$((COUNTER + 1))
done

echo "YDB failed to become healthy after $MAX_WAIT seconds"
exit 1

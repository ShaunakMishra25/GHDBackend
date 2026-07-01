#!/bin/sh
set -e

echo "Running database migrations..."
./server -migrate-only
echo "Migrations complete."

echo "Starting server..."
exec ./server

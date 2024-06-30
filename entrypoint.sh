#!/bin/sh

# Create the .env file with environment variables
echo "TURSO_AUTH_TOKEN=${TURSO_AUTH_TOKEN}" > .env
echo "TURSO_DATABASE_URL=${TURSO_DATABASE_URL}" >> .env
echo "TURSO_JWT_SECRET=${TURSO_JWT_SECRET}" >> .env

# Execute the main process
exec "$@"

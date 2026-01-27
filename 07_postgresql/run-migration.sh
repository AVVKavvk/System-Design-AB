#!/bin/bash
set -e 

PROJECT_ROOT=$(pwd)


DB_URL="${DATABASE_URL:-postgres://myuser:mypassword@localhost:5432/myapp_db?sslmode=disable}"

echo "Waiting for postgres..."
until pg_isready -h localhost -p 5432 -U myuser; do
  echo "Postgres is unavailable - sleeping"
  sleep 2
done

echo "Running migrations..."
# 2. Added -verbose so we can see the table being created
migrate -path "$PROJECT_ROOT/db/users/migrations" -database "$DB_URL" -verbose up

echo "Building application..."
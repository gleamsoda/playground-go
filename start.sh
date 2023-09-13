#!/bin/sh

set -e

echo "run db migrations"
/app/migrate -path /app/migrations -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/playground?parseTime=true" -verbose up

echo "start the app"
exec "$@"
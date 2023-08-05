#!/bin/sh

set -e

echo "run db migrations"
/app/migrate -path /app/migrations -database "mysql://${DB_SOURCE}" -verbose up

echo "start the app"
exec "$@"
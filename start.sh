#!/bin/zsh
set -e

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCEF" -verbose up

echo "start the app"

exec "$@"
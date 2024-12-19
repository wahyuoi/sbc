#!/bin/bash
set -e

# Run migrations
for file in /docker-entrypoint-initdb.d/migrations/*.sql
do
    echo "Running migration: $file"
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE" < "$file"
done

# Run seeds
for file in /docker-entrypoint-initdb.d/seeds/*.sql
do
    echo "Running seed: $file"
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE" < "$file"
done
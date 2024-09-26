#!/bin/zsh

# Function to load .env file
load_env() {
    if [ -f .env ]; then
        export $(grep -v '^#' .env | xargs)
    else
        echo ".env file not found"
        exit 1
    fi
}

# Call the function to load the .env file
load_env

export PGPASSWORD="$DATABASE_PASSWORD"

echo "Syncing Prisma with the Database"
npx prisma migrate dev --skip-generate

pg_dump --host=$DATABASE_HOST --username=$DATABASE_USERNAME --dbname=$DATABASE_NAME --port=$DATABASE_PORT --schema-only |
sed -n '/^CREATE TABLE/,/);$/p' |
sed '/^ALTER TABLE/d' |
awk '
    /^CREATE TABLE public\._prisma_migrations/ {skip=1; next}
    /^CREATE TABLE/ {if (!skip) {print; in_table=1}; next}
    /^);$/ {if (!skip) print; skip=0; in_table=0; next}
    in_table {print}
' > schema.sql

unset PGPASSWORD

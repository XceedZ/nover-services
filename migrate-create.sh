#!/bin/bash

# Script untuk membuat file migrasi baru
# Membutuhkan satu argumen: nama migrasi

if [ -z "$1" ]; then
  echo "❌ Error: Please provide a name for the migration."
  echo "Usage: ./migrate-create.sh <migration_name>"
  exit 1
fi

echo "✨ Creating new migration file: $1"
goose -dir "db/migrations" create $1 sql
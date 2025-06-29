#!/bin/bash

echo "⏪ Reverting last migration..."

# Memuat variabel dari .env dengan aman
set -o allexport
source .env
set +o allexport

# Mengatur variabel yang dibutuhkan goose
export GOOSE_DRIVER=${DB_PROTOCOL}
export GOOSE_DBSTRING="${DB_PROTOCOL}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}&default_query_exec_mode=simple_protocol"

# Menjalankan perintah goose down
goose -dir "db/migrations" down
#!/bin/bash

echo "🚨 PERINGATAN: Perintah ini akan menjalankan semua migrasi 'Down' dan mencoba mengembalikan database ke keadaan kosong."
echo "Tekan Enter untuk melanjutkan, atau Ctrl+C untuk membatalkan."
read

echo "--- Mereset database ---"

# Memuat variabel dari .env dengan aman
set -o allexport
source .env
set +o allexport

# Mengatur variabel yang dibutuhkan goose
export GOOSE_DRIVER=${DB_PROTOCOL}
export GOOSE_DBSTRING="${DB_PROTOCOL}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}&default_query_exec_mode=simple_protocol"

# Menjalankan perintah goose reset dengan mode verbose (-v)
goose -dir "db/migrations" reset -v
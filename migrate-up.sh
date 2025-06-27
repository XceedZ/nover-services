#!/bin/bash

echo "--- LANGKAH 1: MEMUAT FILE .env ---"
if [ ! -f .env ]; then
    echo "❌ FATAL: File .env tidak ditemukan!"
    exit 1
fi
set -o allexport
source .env
set +o allexport
echo "✅ File .env ditemukan dan diproses."

echo ""
echo "--- LANGKAH 2: MEMERIKSA VARIABEL (DEBUGGING) ---"
echo "DRIVER: ${DB_PROTOCOL}"
echo "USER:   ${DB_USER}"
# ... (echo lainnya bisa Anda hapus jika sudah yakin)
if [ -z "${DB_PROTOCOL}" ] || [ -z "${DB_USER}" ] || [ -z "${DB_PASSWORD}" ] || [ -z "${DB_HOST}" ]; then
    echo "❌ FATAL: Variabel database penting tidak ditemukan di .env."
    exit 1
fi
echo "✅ Variabel tampak berhasil di-load ke dalam script."

echo ""
echo "--- LANGKAH 3: MENJALANKAN GOOSE ---"
export GOOSE_DRIVER=${DB_PROTOCOL}

# --- PERUBAHAN UTAMA DI SINI ---
# Menambahkan &default_query_exec_mode=simple_protocol di akhir DSN
export GOOSE_DBSTRING="${DB_PROTOCOL}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}&default_query_exec_mode=simple_protocol"

# Menjalankan goose
goose -dir "db/migrations" up
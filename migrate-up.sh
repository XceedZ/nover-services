#!/bin/bash

# =================================================================
# SCRIPT DEBUGGING UNTUK MENJALANKAN GOOSE
# =================================================================

echo "--- LANGKAH 1: MEMUAT FILE .env ---"

# Cek apakah file .env ada
if [ ! -f .env ]; then
    echo "❌ FATAL: File .env tidak ditemukan di direktori ini!"
    exit 1
fi

# Memuat variabel dari .env dengan cara yang paling aman
set -o allexport
source .env
set +o allexport

echo "✅ File .env ditemukan dan diproses."


echo ""
echo "--- LANGKAH 2: MEMERIKSA VARIABEL (DEBUGGING) ---"

# Kita akan print beberapa variabel untuk memastikan nilainya benar-benar terbaca
# Kita tidak akan print password demi keamanan
echo "DRIVER: ${DB_PROTOCOL}"
echo "USER:   ${DB_USER}"
echo "HOST:   ${DB_HOST}"
echo "PORT:   ${DB_PORT}"
echo "DBNAME: ${DB_NAME}"

# Cek apakah variabel penting kosong atau tidak
if [ -z "${DB_PROTOCOL}" ] || [ -z "${DB_USER}" ] || [ -z "${DB_PASSWORD}" ] || [ -z "${DB_HOST}" ]; then
    echo ""
    echo "❌ FATAL: Satu atau lebih variabel database penting (DB_PROTOCOL, DB_USER, DB_PASSWORD, DB_HOST) tidak ditemukan atau kosong di dalam file .env."
    echo "Pastikan formatnya benar (contoh: KEY=VALUE) tanpa spasi di sekitar =."
    exit 1
fi

echo "✅ Variabel tampak berhasil di-load ke dalam script."


echo ""
echo "--- LANGKAH 3: MENJALANKAN GOOSE ---"

# Mengatur variabel yang dibutuhkan goose secara EKSPLISIT
# Kita tidak lagi bergantung pada goose.conf untuk sementara waktu
export GOOSE_DRIVER=${DB_PROTOCOL}
export GOOSE_DBSTRING="${DB_PROTOCOL}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

# Menjalankan goose. Ia akan otomatis menggunakan variabel GOOSE_DRIVER dan GOOSE_DBSTRING
goose -dir "db/migrations" up
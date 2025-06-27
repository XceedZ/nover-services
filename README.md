# 📘 Nover Services

**Nover Services** adalah backend API service berbasis **Golang** yang digunakan sebagai fondasi backend untuk aplikasi **Nover** — platform membaca novel online.

---

## 🚀 Fitur Utama

* ✅ REST API Backend
* ✅ PostgreSQL Database (Supabase)
* ✅ Database migration dengan **Goose**
* ✅ Modular project structure
* ✅ Support deployment ke cloud / VPS

---

## 📂 Struktur Proyek

```
.
├── cmd/
├── config/
├── db/
│   └── migrations/
├── internal/
├── pkg/
├── .env
├── migrate-create.sh
├── migrate-up.sh
├── migrate-down.sh
├── migrate-reset.sh
├── go.mod
├── go.sum
└── README.md
```

| Folder / File    | Fungsi                              |
| ---------------- | ----------------------------------- |
| `cmd/`           | Entry point aplikasi                |
| `config/`        | Konfigurasi environment             |
| `db/migrations/` | File migrasi database (Goose)       |
| `internal/`      | Business logic utama (modular)      |
| `pkg/`           | Utilities, helpers, middleware, dsb |
| `.env`           | Konfigurasi local environment       |
| `migrate-*.sh`   | Script untuk menjalankan migration  |

---

## ⚙️ Prasyarat

* Go 1.20+
* Supabase (PostgreSQL)
* Goose CLI (v3.x atau terbaru)

---

## 🛠️ Cara Menjalankan Proyek Secara Lokal

### 1. Clone Repository:

```bash
git clone https://github.com/username/nover-services.git
cd nover-services
```

### 2. Siapkan file `.env`

Contoh isi:

```
DB_PROTOCOL=postgres
DB_USER=postgres.your_user
DB_PASSWORD=your_password (URL-encoded jika ada karakter spesial)
DB_HOST=aws-0-ap-southeast-1.pooler.supabase.com
DB_PORT=6543
DB_NAME=postgres
DB_SSL_MODE=require

APP_PORT=8080
APP_ENV=development
```

> ✅ Pastikan `.env` ini tidak di-commit ke git.

---

### 3. Install Dependency:

```bash
go mod tidy
```

---

### 4. Jalankan Migration Database

> Semua script sudah membaca otomatis dari `.env`

#### ✅ Buat migration baru:

```bash
./migrate-create.sh <migration_name>
```

Contoh:

```bash
./migrate-create.sh create_users_table
```

#### ✅ Jalankan migration up:

```bash
./migrate-up.sh
```

#### ✅ Rollback (down 1 step):

```bash
./migrate-down.sh

```
#### ✅ Reset

```bash
./migrate-reset.sh
```

| Script                                 | Fungsi                                                                               |
| -------------------------------------- | ------------------------------------------------------------------------------------ |
| `./migrate-create.sh <migration_name>` | Membuat file migration baru                                                          |
| `./migrate-up.sh`                      | Menjalankan semua migration **up**                                                   |
| `./migrate-down.sh`                    | Rollback migration terakhir (**down 1 step**)                                        |
| `./migrate-reset.sh`                   | **Hati-hati!** Menjalankan semua **down** migration hingga database **kosong total** |

---

## ✅ Menjalankan Server

```bash
go run .
```

Aplikasi akan berjalan di port sesuai `APP_PORT` pada `.env`, contoh:

```
Listening on port 8080
```

---

## 🧱 Tips Supabase + Goose

* Password Supabase biasanya mengandung simbol (`@`, `/`, `:` dll), **pastikan URL Encoded di `.env`!**

* Contoh URL-encode:

  | Karakter | Encode |
  | -------- | ------ |
  | `@`      | `%40`  |
  | `/`      | `%2F`  |
  | `:`      | `%3A`  |

* Pastikan **Goose CLI** kamu versi **3.x atau lebih baru**, agar support **SCRAM-SHA-256 Auth**.

* Kalau error seperti:

```
invalid SCRAM server-final-message
```

> Artinya masalah di **format connection string**, **encoding password**, atau **Goose versi lama**.

---

## ✅ Deployment Notes

* Buat `.env`.
* Pastikan IP server kamu sudah di-whitelist di Supabase.
* Gunakan reverse proxy seperti **Nginx**, **Caddy**, atau **Cloudflare Tunnel** saat deploy.

---

## ✅ Author

| Nama       | Peran              |
| ---------- | ------------------ |
| AlexanderA | Software Developer |

---

## ✅ Lisensi

MIT License
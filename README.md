# CandyzFlorist Backend API (Go + Gin + GORM + MySQL + JWT)

Backend RESTful API untuk sistem manajemen toko bunga **CandyzFlorist**, terintegrasi langsung dengan antarmuka React di folder `candyzflorist-react`.

---

## 🛠️ Tech Stack
- **Bahasa**: Go (v1.20+)
- **Web Framework**: Gin Gonic (`github.com/gin-gonic/gin`)
- **ORM**: GORM (`gorm.io/gorm` + `gorm.io/driver/mysql`)
- **Database**: MySQL
- **Autentikasi**: JWT (`github.com/golang-jwt/jwt/v5`) & Bcrypt (`golang.org/x/crypto/bcrypt`)
- **File Storage**: Local uploads grouped by category slug (`./uploads/<category-slug>/`)

---

## 📂 Struktur Project

```
candyzflorist-go/
├── main.go                  # Entry point aplikasi
├── .env                     # File konfigurasi environment
├── .env.example             # Contoh template environment
├── config/
│   └── database.go          # Koneksi MySQL & Auto Migration
├── models/
│   ├── user.go              # Model User / Admin
│   ├── kategori.go          # Model Kategori
│   ├── produk.go            # Model Katalog Produk
│   ├── produk_image.go      # Model Foto Produk (One-to-Many)
│   └── pengaturan.go        # Model Pengaturan Kontak Toko
├── middleware/
│   ├── auth.go              # JWT Authentication Middleware
│   └── cors.go              # CORS Middleware
├── handlers/
│   ├── auth_handler.go      # Handler Login, Me, Change Password
│   ├── kategori_handler.go  # Handler CRUD Kategori
│   ├── katalog_handler.go   # Handler CRUD Katalog & Upload Foto (Maks 2MB)
│   ├── pengaturan_handler.go# Handler Pengaturan Toko
│   └── dashboard_handler.go # Handler Statistik Dashboard
├── routes/
│   └── routes.go            # Registrasi seluruh API routes
├── seeders/
│   └── seeder.go            # Seeder akun admin awal & kategori
└── helpers/
    ├── response.go          # Standarisasi JSON response
    └── slug.go              # Helper pembuat slug URL otomatis
```

---

## 🚀 Cara Menjalankan

### 1. Buat Database MySQL
Buka MySQL / phpMyAdmin / DBeaver, lalu buat database:
```sql
CREATE DATABASE candyzflorist;
```

### 2. Sesuaikan File `.env`
Buka file `.env` dan sesuaikan kredensial MySQL Anda:
```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=candyzflorist

JWT_SECRET=candyzflorist-super-secret-key-2026
JWT_EXPIRY_HOURS=72

SERVER_PORT=8000
UPLOAD_DIR=./uploads
```

### 3. Jalankan Server
```bash
go run main.go
```
*Tabel database dan data awal (seeder) akan otomatis dibuat saat server pertama kali dijalankan.*

---

## 🔑 Akun Default Seeder

- **Email**: `admin@candyzflorist.com`
- **Password**: `admin123`

---

## 📋 Daftar Endpoint API

### 🌐 Public API (Storefront / Landing Page — Tanpa Login)
| Method | Endpoint | Deskripsi |
|---|---|---|
| **GET** | `/api/public/kategori` | Ambil daftar 11 kategori koleksi katalog |
| **GET** | `/api/public/katalog` | Ambil katalog produk publik (support filter `?kategori=slug` & `?search=keyword`) |
| **GET** | `/api/public/katalog/:identifier` | Ambil detail 1 produk berdasarkan `id` atau `slug` |
| **GET** | `/api/public/pengaturan` | Ambil nomor WhatsApp, medsos, alamat, dan template pesan toko |
| **GET** | `/uploads/*` | Akses publik file gambar buket yang diunggah |

### 🔒 Admin API (Protected — Butuh JWT Bearer Token)
| Method | Endpoint | Deskripsi |
|---|---|---|
| **POST** | `/api/auth/login` | Login admin & dapatkan token JWT (Public) |
| **POST** | `/api/auth/logout` | Logout admin |
| **GET** | `/api/auth/me` | Ambil data profil admin yang sedang login |
| **PUT** | `/api/auth/change-password` | Ganti password admin |
| **GET** | `/api/dashboard/stats` | Ringkasan statistik & 4 produk terbaru |
| **GET** | `/api/kategori` | Ambil daftar semua kategori |
| **GET** | `/api/kategori/:id` | Ambil detail 1 kategori |
| **POST** | `/api/kategori` | Tambah kategori baru (+ auto folder upload) |
| **PUT** | `/api/kategori/:id` | Perbarui kategori |
| **DELETE**| `/api/kategori/:id` | Hapus kategori |
| **GET** | `/api/katalog` | Ambil daftar semua produk katalog |
| **GET** | `/api/katalog/:id` | Ambil detail 1 produk katalog |
| **POST** | `/api/katalog` | Tambah produk (multipart form + multiple images) |
| **PUT** | `/api/katalog/:id` | Edit produk (multipart form + sync images) |
| **DELETE**| `/api/katalog/:id` | Hapus produk & file gambar fisik dari disk |
| **POST** | `/api/katalog/upload` | Standalone upload foto produk |
| **GET** | `/api/pengaturan` | Ambil pengaturan kontak & medsos toko |
| **PUT** | `/api/pengaturan` | Perbarui pengaturan kontak toko |

> **Catatan Header Auth**: Untuk semua endpoint admin (Protected), sertakan header:  
> `Authorization: Bearer <TOKEN_JWT>`

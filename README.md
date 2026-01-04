# Rires Backend

REST API menggunakan **Go 1.25.5**, **Fiber v2**, **GORM**, dan **MySQL**.

## 🚀 Tech Stack

- **Go 1.25.5** - Programming language
- **Fiber v2** - Web framework (Express-like)
- **GORM** - ORM untuk MySQL
- **JWT** - Authentication
- **MySQL** - Database

## 📁 Struktur Project

```
golang-api-tutorial/
├── cmd/api/              # Entry point aplikasi
├── config/               # Konfigurasi & environment
├── internal/
│   ├── controllers/      # HTTP handlers
│   ├── models/          # Database models
│   └── routes/          # Route definitions
└── pkg/
    ├── database/        # Database connection
    ├── middleware/      # Middleware (JWT, logger, etc)
    └── utils/           # Helper functions
```

## ⚙️ Setup

1. Clone repository
2. Copy `.env.example` ke `.env` dan sesuaikan konfigurasi
3. Buat database MySQL dengan nama sesuai di `.env`
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Jalankan aplikasi:
   ```bash
   go run cmd/api/main.go
   ```

## 🔌 API Endpoints

### Authentication
- `POST /api/auth/register` - Register user baru
- `POST /api/auth/login` - Login user

### Users (Protected)
- `GET /api/users/profile` - Get user profile
- `PUT /api/users/profile` - Update user profile

## 📝 License

MIT
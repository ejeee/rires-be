# RIRES Backend API

REST API untuk sistem **Program Kreativitas Mahasiswa (PKM)** Universitas Muhammadiyah Malang menggunakan **Go 1.23**, **Fiber v2**, **GORM**, dan **MySQL**.

## 🚀 Tech Stack

- **Go 1.23+** - Programming language
- **Fiber v2** - Web framework (Express-like)
- **GORM** - ORM untuk MySQL
- **JWT** - Authentication & Authorization
- **MySQL** - Main database + External databases (NEOMAA, NEOMAAREF, SIMPEG)
- **Go Validator** - Request validation

## 📁 Struktur Project

```
rires-be/
├── cmd/api/                    # Entry point aplikasi
├── config/                     # Configuration & environment
├── internal/
│   ├── controllers/            # HTTP handlers
│   ├── dto/
│   │   ├── request/            # Request DTOs
│   │   └── response/           # Response DTOs
│   ├── middleware/             # JWT & role-based middleware
│   ├── models/                 # Database models
│   ├── routes/                 # Route definitions
│   └── services/               # Business logic
└── pkg/
    ├── database/               # Database connections
    ├── services/               # External services
    └── utils/                  # Helper functions
```

## ⚙️ Setup & Installation

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd rires-be
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment**
   ```bash
   cp .env.example .env
   ```

4. **Run application**
   ```bash
   go run cmd/api/main.go
   ```

## 🔌 API Endpoints (27+)

See detailed documentation in the codebase or API docs.

## 📜 License

MIT License

---
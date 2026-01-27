# RIRES Backend API

REST API untuk sistem **Program Kreativitas Mahasiswa (PKM)** Universitas Muhammadiyah Malang. Backend ini dibangun menggunakan **Go 1.23**, **Fiber v2**, dan **GORM**.

## 🚀 Tech Stack

- **Go 1.23+** - Main programming language.
- **Fiber v2** - Fast and minimal web framework.
- **GORM** - ORM for database interaction.
- **JWT** - Secure authentication and role-based access control.
- **MySQL** - Integrated with multiple databases:
  - `Main DB`: Core tables for PKM management.
  - `NEOMAA`: Student data integration.
  - `NEOMAAREF`: Reference data (Fakultas, Prodi).
  - `SIMPEG`: Employee/Reviewer data integration.
- **Go Validator** - Robust request validation.

## 📁 Project Structure

```text
rires-be/
├── cmd/
│   └── api/                # Application entry point (main.go)
├── config/                 # Configuration management & environment loading
├── docs/                   # Swagger documentation & API specs
├── internal/
│   ├── controllers/        # HTTP handlers (logic for each route)
│   ├── dto/                # Data Transfer Objects
│   │   ├── request/        # Request body structures & validation rules
│   │   └── response/       # Standardized API response structures
│   ├── middleware/         # Security middlewares (JWT, Role checking)
│   ├── models/             # GORM models for local database
│   │   └── external/       # Models for external database integrations
│   └── routes/             # Central route setup & group definitions
└── pkg/
    ├── database/           # Multi-database connection setup
    ├── services/           # Business logic & external data integration
    └── utils/              # Common helpers (JWT, Response, Strings)
```

## 🛠️ Getting Started

1.  **Clone the Repository**
    ```bash
    git clone <repository-url>
    cd rires-be
    ```

2.  **Install Dependencies**
    ```bash
    go mod download
    ```

3.  **Environment Configuration**
    Copy `.env.example` to `.env` and fill in necessary database credentials.
    ```bash
    cp .env.example .env
    ```

4.  **Run Development Server**
    ```bash
    go run cmd/api/main.go
    ```

## 🔑 Key Features

- **Multi-Role Authentication**: Support for Admin, Mahasiswa, and Reviewer logins.
- **Reviewer Assignment**: Automated and manual plotting of reviewers for PKM titles and proposals.
- **Flexible Review Flow**: Support for revision, acceptance, and rejection cycles.
- **Database Integration**: Seamless synchronization with UMM's internal systems (SIMPEG, NEOMAA).

## 📜 License

This project is licensed under the MIT License.

---

tes perubahan
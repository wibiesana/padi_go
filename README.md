<div align="center">

# 🌾 PADI REST API FRAMEWORK (GO)

**The Industrial-Grade, Zero-Bloat REST API Engine for Go**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![Architecture](https://img.shields.io/badge/Architecture-Code_Generation-success?style=flat-square)](#)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](#)

> **"Stop wrestling with massive frameworks. Start shipping APIs in seconds."**

</div>

Welcome to the **Padi REST API Go Framework** — a high-performance Go port of the Padi REST API framework designed for developers who crave raw speed, absolute clarity, and an uncompromising Developer Experience (DX).

---

## ✨ Features (100% Feature Parity with Padi PHP)

- ⚡ **Zero-Bloat & Raw Speed**: Minimal latency HTTP pipeline powered by standard `database/sql`, `chi` router, and pure Go architecture (no external ORM dependencies).
- 🌾 **Pure Native ActiveRecord & Fluent Query Builder**: Built-in `model.Save()`, `model.Delete()`, `(User{}).Find(id)`, `query.New("table").Where(...).Paginate(...)` persis seperti Padi PHP.
- 🛠️ **Signature Auto CRUD Generator**: Generate Base Models, Custom Models, Controllers, and Routes in seconds from DB schema with `padi g <table_name>` or batch all tables with `padi ga`.
- 🔐 **Security Built-In**: JWT Authentication, Bcrypt hashing, Request Rate Limiting, and CORS management.
- 🗄️ **Database Agnostic**: Native support for **SQLite (pure Go)**, **MySQL**, and **PostgreSQL** via standard `database/sql` connection pooling.
- 🚀 **Built-in Migrations**: Automatic migration registry and rollback system (`padi migrate`, `padi migrate:rollback`).
- 📦 **L1 & L2 Caching**: In-Memory + Redis / File caching with TTL, auto-eviction, and `cache.Remember()`.
- 📁 **File & Media Storage**: Secure file uploads with extension blacklists/whitelists, MIME validation, and automatic URL resolution (`storage.SaveUploadedFile()`, `storage.URL()`).
- 📧 **Email Engine**: SMTP mailer with TLS/SSL support (`email.Send()`).
- 👷 **Async Background Queue & Workers**: Database-backed asynchronous queue and worker runner (`queue.Push()`, `padi queue:work`).
- 📡 **Real-time Pub/Sub**: Native Server-Sent Events (SSE) broadcaster (`realtime.Publish()`, `realtime.SubscribeSSE()`).
- ✅ **Type-Safe Validation**: Field-level request validation with clean JSON error responses.

---

## 📁 Project Structure

```text
padi_rest_api_go_framework/
├── app/
│   ├── Controllers/            # HTTP Handlers (AuthController, UserController, etc.)
│   ├── Models/                 # Custom domain models (business logic, ActiveRecord methods)
│   │   └── Base/               # Auto-generated base models (DO NOT EDIT directly)
│   └── Routes/
│       └── api.go              # API route registry with versioning
├── cmd/
│   └── padi/                   # Padi CLI tool
│       └── main.go
├── config/                     # Configuration and Environment binding
├── database/
│   ├── migrations/             # Database migrations
│   └── database.sqlite         # SQLite default database file
├── pkg/
│   └── core/                   # Framework Core (Zero External ORM)
│       ├── auth/               # JWT & Password hashing
│       ├── cache/              # Memory, Redis, and File Caching
│       ├── config/             # .env loader & config
│       ├── database/           # Native database/sql pool manager (SQLite, MySQL, Postgres)
│       ├── email/              # SMTP mailer
│       ├── generator/          # Auto-CRUD & scaffolding engine
│       ├── middleware/         # AuthRequired, CORS, Logger, RateLimiter, Recoverer
│       ├── migrator/           # Schema migration runner
│       ├── model/              # Native Padi ActiveRecord engine (Find, Save, Delete)
│       ├── query/              # Native Fluent Query Builder & pagination
│       ├── queue/              # Background Job Queue and Worker
│       ├── realtime/           # SSE Realtime Pub/Sub Hub
│       ├── response/           # Standardized JSON response formatting
│       ├── router/             # Router wrapper & param extractors
│       ├── storage/            # File upload & static URL helpers
│       └── validator/          # Struct validation and JSON binding
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── main.go                     # Server application entrypoint
└── README.md
```

---

## 🛠️ Padi CLI Commands

| Command | Alias | Description |
| :--- | :--- | :--- |
| `go run ./cmd/padi serve` | | Start the HTTP API server |
| `go run ./cmd/padi migrate` | | Run pending database migrations |
| `go run ./cmd/padi migrate:rollback` | | Rollback last batch of migrations |
| `go run ./cmd/padi generate:crud <table_name>` | `padi g <table_name>` | Auto-generate Base Model, Custom Model, Controller, and Route snippet |
| `go run ./cmd/padi generate:crud-all` | `padi ga` | Auto-generate CRUD for **ALL** database tables |
| `go run ./cmd/padi update` | | Update Padi Core dependencies and modules to latest versions |
| `go run ./cmd/padi queue:work [queue_name]` | `padi queue` | Run the background queue worker |
| `go run ./cmd/padi make:model <name>` | `padi create:model` | Create a new model file |
| `go run ./cmd/padi make:controller <name>` | `padi create:controller` | Create a new controller file |
| `go run ./cmd/padi make:migration <name>` | `padi create:migration` | Create a new migration file |

---

## 📡 Standard API Response Format

```json
{
  "status": 200,
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "meta": {
    "total": 100,
    "per_page": 15,
    "current_page": 1,
    "last_page": 7,
    "from": 1,
    "to": 15
  }
}
```

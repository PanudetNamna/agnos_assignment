# Agnos Backend – Hospital Middleware

A Hospital Middleware API built with Go, Gin, PostgreSQL, Docker, and Nginx.  
Acts as a middleware layer between client applications and Hospital Information Systems (HIS).

## Tech Stack

- **Go** + **Gin Framework**
- **PostgreSQL** (via GORM)
- **Docker** + **Docker Compose**
- **Nginx** (reverse proxy)
- **JWT** authentication

## Project Structure

```
.
├── cmd/
│   └── main.go
├── config/
│   ├── config.yaml
│   └── secret.env
├── internal/
│   ├── adapter/
│   │   └── his/
│   │       └── client.go         # HIS client implementation
│   ├── config/
│   │   ├── common.go             # Config loader
│   │   └── model.go              # Config struct
│   ├── handlers/
│   │   └── http/
│   │       ├── health_check/     # Health check handler + test
│   │       ├── patient/          # Patient handler
│   │       ├── staff/            # Staff handler
│   │       └── router.go
│   ├── middleware/
│   │   └── auth.go               # JWT auth middleware
│   ├── models/
│   │   ├── his.go
│   │   ├── hospital.go
│   │   ├── patient.go
│   │   └── staff.go
│   ├── port/
│   │   ├── handlers.go           # Handler interface
│   │   ├── his_client.go         # HIS client interface
│   │   ├── repositories.go       # Repository interfaces
│   │   └── services.go           # Service interfaces
│   ├── repositories/
│   │   ├── hospital/
│   │   ├── patient/
│   │   └── staff/
│   ├── services/
│   │   ├── patient/
│   │   └── staff/
│   └── utility/
│       ├── response.go           # HTTP response helpers
│       └── token.go              # JWT helpers
├── nginx/
│   └── nginx.conf
├── Dockerfile
└── docker-compose.yml
```

## Architecture

This project follows **Hexagonal Architecture** (Ports & Adapters).

```
Handler → Service → Repository → PostgreSQL
                 ↘ HIS Client → External HIS API
```

- **port/** defines all interfaces
- **adapter/** contains external integrations (HIS)
- **services/** contains business logic, depends only on interfaces
- **handlers/** contains HTTP layer, depends only on service interfaces

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /health | No | Health check |
| POST | /staff/create | No | Create a new staff member |
| POST | /staff/login | No | Login and get JWT token |
| POST | /patient/search | JWT | Search patients from HIS + DB |

---

### GET /health

Response `200`:
```json
{ "message": "ok" }
```

---

### POST /staff/create

Request:
```json
{
  "username": "john",
  "password": "secret",
  "hospital": "Hospital A"
}
```

Response `201`:
```json
{
  "message": "success",
  "data": {
    "staff_id": "uuid",
    "hospital_id": "uuid"
  }
}
```

---

### POST /staff/login

Request:
```json
{
  "username": "john",
  "password": "secret",
  "hospital": "Hospital A"
}
```

Response `200`:
```json
{
  "message": "success",
  "data": {
    "token": "<JWT>"
  }
}
```

---

### POST /patient/search

Headers: `Authorization: Bearer <JWT>`

Request body (all optional):
```json
{
  "national_id": "1234567890123",
  "passport_id": "",
  "first_name": "John",
  "middle_name": "",
  "last_name": "Doe",
  "date_of_birth": "1990-01-01",
  "phone_number": "",
  "email": ""
}
```

Response `200`:
```json
{
  "message": "success",
  "data": {
    "count": 1,
    "patients": [...]
  }
}
```

> If `national_id` or `passport_id` is provided, the system will fetch and sync patient data from the HIS before querying the local database.

---

## Error Response Format

```json
{
  "message": "error description"
}
```

## Run with Docker

```bash
docker compose up --build
```

App will be available at `http://localhost` (via Nginx → port 80).

## Run Tests

```bash
go test ./...
```

## Configuration

| File | Description |
|------|-------------|
| `config/config.yaml` | App config (DB, server, HIS) |
| `config/secret.env` | Secrets (JWT secret key) |


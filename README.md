# Ainyx User API

A high-performance RESTful API for user management built with Go, Fiber, and PostgreSQL. This API provides full CRUD operations for user data with features like pagination, validation, structured logging, and age calculation.

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-316192?style=flat&logo=postgresql)

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Development](#development)
- [Docker Deployment](#docker-deployment)
- [Testing](#testing)
- [Contributing](#contributing)

---

## Features

- **Full CRUD Operations**: Create, Read, Update, and Delete users
- **Automatic Age Calculation**: Computes user age from date of birth
- **Pagination Support**: Efficiently list users with customizable page sizes
- **Input Validation**: Robust validation using `go-playground/validator`
- **Structured Logging**: Production-ready logging with `uber-go/zap`
- **Type-Safe SQL**: Generated SQL queries using `sqlc`
- **Fast HTTP Framework**: High-performance routing with Fiber v2
- **Docker Support**: Multi-stage Docker build for optimized images
- **Clean Architecture**: Separation of concerns with handler, service, and repository layers

---

## Tech Stack

| Technology | Purpose |
|------------|---------|
| [Go 1.23+](https://golang.org/) | Programming Language |
| [Fiber v2](https://gofiber.io/) | Web Framework |
| [PostgreSQL](https://www.postgresql.org/) | Database |
| [pgx/v5](https://github.com/jackc/pgx) | PostgreSQL Driver |
| [sqlc](https://sqlc.dev/) | Type-safe SQL Generation |
| [Zap](https://github.com/uber-go/zap) | Structured Logging |
| [Validator](https://github.com/go-playground/validator) | Input Validation |
| [godotenv](https://github.com/joho/godotenv) | Environment Configuration |

---

## Project Structure

```
ainyx-user-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── config/
│   └── config.go                # Configuration management
├── db/
│   ├── migrations/
│   │   └── 001_create_users.sql # Database migrations
│   └── sqlc/
│       ├── queries/
│       │   └── users.sql        # SQL queries for sqlc
│       ├── db.go                # Generated database interface
│       ├── models.go            # Generated models
│       └── users.sql.go         # Generated query functions
├── internal/
│   ├── handler/
│   │   └── user_handler.go      # HTTP request handlers
│   ├── logger/
│   │   └── logger.go            # Zap logger configuration
│   ├── middleware/              # HTTP middleware
│   ├── models/
│   │   ├── user.go              # Request/Response DTOs
│   │   ├── user_test.go         # Model tests
│   │   └── validator.go         # Validation utilities
│   ├── repository/
│   │   └── user_repository.go   # Data access layer
│   ├── routes/
│   │   └── routes.go            # Route definitions
│   └── service/
│       └── user_service.go      # Business logic layer
├── .env                         # Environment variables (not in git)
├── .gitignore                   # Git ignore rules
├── Dockerfile                   # Docker multi-stage build
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── sqlc.yaml                    # sqlc configuration
└── README.md                    # This file
```

---

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.23+**: [Download Go](https://golang.org/dl/)
- **PostgreSQL 15+**: [Download PostgreSQL](https://www.postgresql.org/download/)
- **sqlc** (optional, for regenerating queries): [Install sqlc](https://sqlc.dev/)
- **Docker** (optional, for containerized deployment): [Install Docker](https://docs.docker.com/get-docker/)

---

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/muharib-0/ainyx-user-api.git
cd ainyx-user-api
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up the Database

Create a PostgreSQL database and run the migration:

```sql
-- Connect to PostgreSQL and create the database
CREATE DATABASE ainyx_users;

-- Connect to the new database and run migration
\c ainyx_users

-- Create the users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    dob DATE NOT NULL
);
```

Or run the migration file directly:

```bash
psql -U postgres -d ainyx_users -f db/migrations/001_create_users.sql
```

---

## Configuration

### Environment Variables

Create a `.env` file in the project root with the following variables:

```env
# Server Configuration
SERVER_PORT=3000

# Database Configuration
DATABASE_URL=postgres://username:password@localhost:5432/ainyx_users?sslmode=disable
DB_DRIVER=postgres
```

### Configuration Options

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Port the server listens on | `3000` |
| `DATABASE_URL` | PostgreSQL connection string | (required) |
| `DB_DRIVER` | Database driver | `postgres` |

---

## Running the Application

### Development Mode

```bash
go run ./cmd/server
```

### Build and Run

```bash
# Build the binary
go build -o server ./cmd/server

# Run the binary
./server
```

### Verify the Server is Running

```bash
curl http://localhost:3000/health
```

Expected response:

```json
{
  "status": "ok",
  "message": "Server is running"
}
```

---

## API Documentation

### Base URL

```
http://localhost:3000/api/v1
```

### Endpoints

#### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Check server health |

#### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/users` | Create a new user |
| GET | `/api/v1/users/:id` | Get a user by ID |
| GET | `/api/v1/users` | List all users (paginated) |
| PUT | `/api/v1/users/:id` | Update a user |
| DELETE | `/api/v1/users/:id` | Delete a user |

---

### Request & Response Examples

#### Create User

**Request:**

```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "dob": "1990-05-15"
  }'
```

**Response (201 Created):**

```json
{
  "id": 1,
  "name": "John Doe",
  "dob": "1990-05-15"
}
```

#### Get User by ID

**Request:**

```bash
curl http://localhost:3000/api/v1/users/1
```

**Response (200 OK):**

```json
{
  "id": 1,
  "name": "John Doe",
  "dob": "1990-05-15",
  "age": 34
}
```

#### List Users (Paginated)

**Request:**

```bash
curl "http://localhost:3000/api/v1/users?page=1&page_size=10"
```

**Response (200 OK):**

```json
{
  "users": [
    {
      "id": 1,
      "name": "John Doe",
      "dob": "1990-05-15",
      "age": 34
    },
    {
      "id": 2,
      "name": "Jane Smith",
      "dob": "1985-08-22",
      "age": 39
    }
  ],
  "total": 2,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

#### Update User

**Request:**

```bash
curl -X PUT http://localhost:3000/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "dob": "1990-05-15"
  }'
```

**Response (200 OK):**

```json
{
  "id": 1,
  "name": "John Updated",
  "dob": "1990-05-15"
}
```

#### Delete User

**Request:**

```bash
curl -X DELETE http://localhost:3000/api/v1/users/1
```

**Response (204 No Content)**

---

### Validation Rules

| Field | Rules |
|-------|-------|
| `name` | Required, 1-255 characters |
| `dob` | Required, format: `YYYY-MM-DD` |

### Pagination Parameters

| Parameter | Description | Default | Max |
|-----------|-------------|---------|-----|
| `page` | Page number | 1 | - |
| `page_size` | Items per page | 10 | 100 |

---

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    dob DATE NOT NULL
);
```

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-incrementing ID |
| `name` | TEXT | NOT NULL | User's full name |
| `dob` | DATE | NOT NULL | Date of birth |

---

## Development

### Regenerating SQLC Code

If you modify the SQL queries, regenerate the Go code:

```bash
sqlc generate
```

### Running Tests

```bash
go test ./...
```

### Running Tests with Coverage

```bash
go test -cover ./...
```

### Linting

```bash
go vet ./...
```

---

## Docker Deployment

### Build Docker Image

```bash
docker build -t ainyx-user-api .
```

### Run Container

```bash
docker run -d \
  --name ainyx-api \
  -p 3000:3000 \
  -e SERVER_PORT=3000 \
  -e DATABASE_URL="postgres://user:password@host.docker.internal:5432/ainyx_users?sslmode=disable" \
  ainyx-user-api
```

### Docker Compose (Optional)

Create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "3000:3000"
    environment:
      - SERVER_PORT=3000
      - DATABASE_URL=postgres://postgres:postgres@db:5432/ainyx_users?sslmode=disable
    depends_on:
      - db

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=ainyx_users
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./db/migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"

volumes:
  postgres_data:
```

Run with Docker Compose:

```bash
docker-compose up -d
```

---

## Testing

### Unit Tests

```bash
go test ./internal/models/...
```

### Integration Tests

```bash
go test ./... -tags=integration
```

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---



**Muharib** - [@muharib-0](https://github.com/muharib-0)

Project Link: [https://github.com/muharib-0/ainyx-user-api](https://github.com/muharib-0/ainyx-user-api)

---



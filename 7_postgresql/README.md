# PostgreSQL User Management API

A RESTful API built with Go (Echo framework) for managing user data with PostgreSQL database. This project demonstrates clean architecture principles with proper separation of concerns across API handlers, service layer, and data models.

## Why Golang Migrate?

**Golang Migrate** is a powerful database migration tool that helps manage database schema changes in a version-controlled and reversible manner. Here's why it's essential for this project:

### Benefits:

- **Version Control**: Track all database schema changes alongside your code
- **Reversibility**: Roll back changes safely with `.down.sql` migrations
- **Team Collaboration**: Ensures all team members have the same database schema
- **Environment Consistency**: Maintain identical schemas across development, staging, and production
- **CI/CD Integration**: Automate database migrations in deployment pipelines
- **Database Agnostic**: Supports multiple databases (PostgreSQL, MySQL, SQLite, etc.)
- **No ORM Lock-in**: Works independently of your ORM choice or raw SQL usage

### Key Features:

- Sequential or timestamp-based migration ordering
- Support for raw SQL migrations
- CLI tool for easy migration management
- Programmatic API for embedding in applications
- Dirty state detection and recovery

## Installation & Setup

### Prerequisites

- Go 1.19 or higher
- PostgreSQL 12 or higher
- Docker & Docker Compose (optional)

### Install Golang Migrate

```bash
# Install migrate CLI with PostgreSQL support
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Verify installation:

```bash
migrate -version
```

![](./images/image.png)

## Using Golang Migrate

### Creating Migrations

Create a new migration file:

```bash
migrate create -ext sql -dir db/users/migrations -seq <migration_name>
```

Example:

```bash
migrate create -ext sql -dir db/users/migrations -seq users
```

This creates two files:

- `000001_users.up.sql` - Schema changes to apply
- `000001_users.down.sql` - Schema changes to rollback

### Running Migrations

**Apply all pending migrations:**

```bash
migrate -path db/users/migrations -database "postgresql://username:password@localhost:5432/dbname?sslmode=disable" up
```

**Rollback last migration:**

```bash
migrate -path db/users/migrations -database "postgresql://username:password@localhost:5432/dbname?sslmode=disable" down 1
```

**Rollback all migrations:**

```bash
migrate -path db/users/migrations -database "postgresql://username:password@localhost:5432/dbname?sslmode=disable" down
```

**Check migration version:**

```bash
migrate -path db/users/migrations -database "postgresql://username:password@localhost:5432/dbname?sslmode=disable" version
```

**Force migration version (use carefully):**

```bash
migrate -path db/users/migrations -database "postgresql://username:password@localhost:5432/dbname?sslmode=disable" force <version>
```

### Using the Migration Script

This project includes `run-migration.sh` for convenient migration management:

```bash
chmod +x run-migration.sh
./run-migration.sh
```

## Official Documentation

- **Golang Migrate GitHub**: https://github.com/golang-migrate/migrate
- **Documentation**: https://github.com/golang-migrate/migrate/tree/master/database/postgres
- **CLI Usage**: https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md

## API Endpoints

### Base URL

```
http://localhost:8080
```

### User Management

#### 1. Create User

- **Endpoint**: `POST /users`
- **Description**: Create a new user in the database
- **Request Body**:

```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "name": "John Doe",
  "age": 25
}
```

- **Success Response**: `201 Created`

```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "name": "John Doe",
  "age": 25,
  "created_at": "2026-01-19T10:30:00Z"
}
```

- **Error Response**: `400 Bad Request` | `500 Internal Server Error`

---

#### 2. Get All Users

- **Endpoint**: `GET /users`
- **Description**: Fetch all users from the database
- **Success Response**: `200 OK`

```json
[
  {
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe",
    "name": "John Doe",
    "age": 25,
    "created_at": "2026-01-19T10:30:00Z"
  },
  {
    "id": 2,
    "email": "jane@example.com",
    "username": "janedoe",
    "name": "Jane Doe",
    "age": 28,
    "created_at": "2026-01-19T11:00:00Z"
  }
]
```

---

#### 3. Get User by ID

- **Endpoint**: `GET /users/{id}`
- **Description**: Fetch a specific user by their ID
- **Path Parameters**:
  - `id` (string, required): User ID
- **Success Response**: `200 OK`

```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "name": "John Doe",
  "age": 25,
  "created_at": "2026-01-19T10:30:00Z"
}
```

- **Error Response**: `404 Not Found`

---

#### 4. Update User

- **Endpoint**: `PUT /users/{id}`
- **Description**: Update an existing user's information
- **Path Parameters**:
  - `id` (string, required): User ID
- **Request Body**:

```json
{
  "email": "newemail@example.com",
  "username": "johndoe_updated",
  "name": "John Doe Jr.",
  "age": 26
}
```

- **Success Response**: `200 OK`

```json
{
  "id": 1,
  "email": "newemail@example.com",
  "username": "johndoe_updated",
  "name": "John Doe Jr.",
  "age": 26,
  "created_at": "2026-01-19T10:30:00Z"
}
```

- **Error Response**: `400 Bad Request` | `404 Not Found`

---

#### 5. Delete User

- **Endpoint**: `DELETE /users/{id}`
- **Description**: Delete a user by their ID
- **Path Parameters**:
  - `id` (string, required): User ID
- **Success Response**: `204 No Content`
- **Error Response**: `400 Bad Request` | `404 Not Found`

---

## Project Structure

```
.
├── api/              # HTTP handlers and routes
├── config/           # Configuration files
├── db/               # Database migrations and connection
│   └── users/
│       └── migrations/
├── docs/             # Swagger/API documentation
├── models/           # Data models and structs
├── service/          # Business logic layer
├── main.go           # Application entry point
├── docker-compose.yml
├── run-migration.sh  # Migration helper script
└── README.md
```

## Running the Application

### With Docker Compose

```bash
docker-compose up -d
```

### Locally

```bash
# Run migrations
./run-migration.sh

# Start the application
go run main.go
```

## API Documentation

Once the application is running, access the Swagger documentation at:

```
http://localhost:8080/docs/index.html
```

## Technologies Used

- **Go** - Programming language
- **Echo** - Web framework
- **PostgreSQL** - Database
- **Golang Migrate** - Database migrations
- **Swagger** - API documentation
- **Docker** - Containerization

## PostgreSQL Notes

https://github.com/AVVKavvk/System-Design-AB/tree/main/7_postgresql/PostGres.pdf

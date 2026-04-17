# Task Tracker API

REST API for task management with user authentication.

This project was built as a backend portfolio project using Go. It demonstrates:

- RESTful API design
- JWT authentication
- PostgreSQL integration
- SQL migrations
- Dockerized local setup
- Protected routes with middleware
- Clean project structure

---

## Tech Stack

- Go 1.25
- Chi v5
- PostgreSQL 17
- pgx / pgxpool
- JWT authentication
- bcrypt password hashing
- Docker + Docker Compose
- golang-migrate for database migrations

---

## Run Project

### Requirements

- Docker
- Docker Compose

### Start Application

```bash
docker compose up --build
```

This command will:

1. Start PostgreSQL container
2. Wait until database becomes healthy
3. Run migrations automatically
4. Build Go application image
5. Start API server on port `8080`

API will be available at:

```text
http://localhost:8080
```

---

## Database Migrations

Migrations are executed automatically by the `migrate` container during startup.

Current migrations:

- create `tasks` table
- create `users` table
- add `user_id` relation to tasks

If migrations fail, application container will not start.

---

## API Endpoints

## Health Check

### GET `/health`

```bash
curl http://localhost:8080/health
```

Response:

```json
{
  "status": "ok"
}
```

---

## Authentication

## Register

### POST `/auth/register`

```bash
curl -X POST http://localhost:8080/auth/register \
-H "Content-Type: application/json" \
-d '{
  "email":"test@example.com",
  "password":"secret123"
}'
```

Response:

```text
201 Created
```

---

## Login

### POST `/auth/login`

```bash
curl -X POST http://localhost:8080/auth/login \
-H "Content-Type: application/json" \
-d '{
  "email":"test@example.com",
  "password":"secret123"
}'
```

Response:

```json
{
  "token":"YOUR_JWT_TOKEN"
}
```

---

## Tasks (Protected Routes)

All `/tasks` endpoints require:

```text
Authorization: Bearer YOUR_JWT_TOKEN
```

---

## Get Tasks

### GET `/tasks`

```bash
curl http://localhost:8080/tasks \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Optional query params:

- `status=new`
- `status=in_progress`
- `status=done`
- `limit=10`

---

## Create Task

### POST `/tasks`

```bash
curl -X POST http://localhost:8080/tasks \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-H "Content-Type: application/json" \
-d '{
  "name":"Buy groceries"
}'
```

---

## Get Task By ID

### GET `/tasks/{taskID}`

```bash
curl http://localhost:8080/tasks/TASK_ID \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Update Task

### PUT `/tasks/{taskID}`

```bash
curl -X PUT http://localhost:8080/tasks/TASK_ID \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-H "Content-Type: application/json" \
-d '{
  "name":"Updated task",
  "status":"done"
}'
```

---

## Delete Task

### DELETE `/tasks/{taskID}`

```bash
curl -X DELETE http://localhost:8080/tasks/TASK_ID \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Task Status Values

Allowed statuses:

- `new`
- `in_progress`
- `done`

---

## Project Structure

```text
cmd/api              application entrypoint
internal/auth        JWT logic
internal/config      config loader
internal/handlers    HTTP handlers
internal/middleware  auth middleware
internal/models      DTO / models
internal/storage     database layer
migrations           SQL migrations
```

---

## Notes

This project is focused on backend fundamentals:

- authentication
- route protection
- clean handlers
- repository pattern
- SQL work
- Dockerized development flow


# Task Management System (Go, Gin, GORM, MySQL)

## Problem Breakdown and Design Decisions

### Objective
Build a small Task microservice that exposes REST APIs for CRUD, supports pagination and status filtering, and is easy to run locally and extend.

### Functional Scope
- Create, read (single + list), update, delete **Task**.
- List supports `page`, `page_size`, and `status` filter (`Pending|InProgress|Completed`).
- Health endpoints: `/healthz`, `/ready`.

### Architecture (Single Responsibility)
- **transport/http** — routing, request binding/validation, HTTP status codes, uniform response envelope.
- **service** — business rules (required title, allowed status, pagination caps), orchestration.
- **repository** — GORM-based persistence behind an interface; DB-agnostic at the boundary.
- **models** — domain types (`Task`, `Status`) without framework concerns.

### Data Model
- `Task { id, title, description, status, dueDate, createdAt, updatedAt }`.
- `status`: `VARCHAR(20)` with default `Pending` (portable; can switch to MySQL `ENUM`).
- Indexes: `created_at` (listing order), optional `status` (filtered queries).

### API Design
- Resource-oriented routes under `/tasks` with standard verbs:
  - `POST /tasks`, `GET /tasks`, `GET /tasks/{id}`, `PUT /tasks/{id}`, `DELETE /tasks/{id}`.
- Consistent JSON envelope for success/error; meaningful HTTP codes (201/200/204/400/404/500).
- Pagination response includes `meta { page, pageSize, total, totalPages }`.

### Validation & Errors
- Validate JSON body (title required, status must be one of the enum values).
- Centralized error helper returns envelope:  
  `{"status":"error","statusCode":..., "error": {"code":"...", "message":"..."}, "requestId":"..."}`.

### Response Envelope
- All responses include: `status`, `statusCode`, optional `data`, optional `meta`, `requestId`.
- Predictable shape simplifies clients and log correlation.

### Database & Migrations
- MySQL 8.0 (via Docker). GORM automigration on startup creates/updates `tasks` table.
- Connection pool tuned (open/idle limits, lifetimes).

### Configuration
- Sensible defaults (no env required): host `127.0.0.1`, port `3307`, user `appuser`, pass `apppass`, db `tasksdb`.
- Optional single `DATABASE_DSN` override or `DB_*` parts.

### Observability & Ops
- `X-Request-ID` middleware for correlation.
- `/healthz` and `/ready` for liveness/readiness.
- Graceful shutdown: drain server and **close DB pool** on SIGINT/SIGTERM.

### Scalability
- Stateless application → horizontal scale behind a load balancer.
- DB: add read replicas/caching as load grows; keep queries simple and indexed.

### Testing Approach
- Postman collection for quick test.

### Trade-offs
- GORM chosen for speed of delivery over raw SQL; easy migrations but less control.
- Uniform envelope adds a small payload overhead, but simplifies client handling and troubleshooting.


## Instructions to Run the Service

### Prerequisites
- Go **1.22+**
- Docker Desktop (WSL2 on Windows is fine)

---

### 1) Start MySQL (Docker)
> Compose maps **host 3307 → container 3306** by default. If 3307 is busy, change the mapping in `docker-compose.yml`.

```bash
# from repo root
docker compose up -d mysql
docker compose ps mysql   # check Ports (e.g., 0.0.0.0:3307->3306/tcp) and Health
```
### 2) Run the golang server
```
go run main.go
```

Test these api with the Postman or any other tool
```bash
curl -X GET "http://localhost:8080/healthz"
curl -X GET "http://localhost:8080/ready"
curl -X POST "http://localhost:8080/tasks" -H "Content-Type: application/json" -d '{"title":"My first task","description":"created from Postman","status":"Pending"}'
curl -X GET "http://localhost:8080/tasks?page=1&page_size=10"
curl -X GET "http://localhost:8080/tasks?status=Pending&page=1&page_size=10"
curl -X GET "http://localhost:8080/tasks/<TASK_ID>"
curl -X PUT "http://localhost:8080/tasks/<TASK_ID>" -H "Content-Type: application/json" -d '{"title":"My first task (updated)","status":"InProgress"}'
curl -X DELETE "http://localhost:8080/tasks/<TASK_ID>"
```
### Response Envelope (applies to all endpoints)

```json
{
  "status": "success | error",
  "statusCode": 200,
  "data": {},
  "error": { "code": "BAD_REQUEST", "message": "..." },
  "meta": {},
  "requestId": "..."
} 
```

## How This Service Demonstrates Microservices Concepts

### 1) Single Responsibility & Clear Boundaries
- **Task Service** owns only the **Task** domain (validation, lifecycle, persistence).
- Layers follow SRP:
  - `transport/http`: routing, binding, HTTP codes, uniform response envelope.
  - `service`: business rules, orchestration, pagination limits.
  - `repository`: data access via interface (GORM/MySQL).
  - `models`: domain types (`Task`, `Status`) without framework coupling.
- No cross-service joins; external data (e.g., users) referenced by `userId` only.

### 2) API Design & Contract
- Resource-oriented REST (`/tasks`) with standard verbs and status codes.
- Consistent JSON **envelope** (`status`, `statusCode`, `data|error`, `meta`, `requestId`) for predictable clients.
- Health/readiness endpoints (`/healthz`, `/ready`) for deploys and automation.
- Backward-compatibility path via `/v1` versioning when needed.

### 3) Scalability & Resilience
- **Stateless** application → horizontal scaling behind a load balancer/API gateway.
- DB connection pooling, simple indexed queries (`status`, `created_at`), and pagination for controlled load.
- Read replicas and optional caching (e.g., Redis) can be added without changing the API.
- Graceful shutdown drains in-flight requests; timeouts/retries can be applied to outbound calls.

### 4) Data Ownership & Independence
- Each service owns its database schema (here: `tasks`), enabling independent deploys and schema evolution.
- Repository interface makes it trivial to swap storage (MySQL → Postgres) or add a cache layer.

### 5) Inter-Service Communication Options
- **REST (JSON)**: simple synchronous lookups (e.g., Task Service → User Service `GET /users/{id}`) with request ID propagation.
- **gRPC**: high-throughput, strongly typed internal calls; deadlines, retries, and mTLS via a service mesh.
- **Events (Kafka/NATS/RabbitMQ)**: async workflows and fan-out (`task.created`, `task.completed`), using the **Outbox Pattern** for reliable delivery; idempotent consumers and a schema registry.

### 6) Observability & Operability
- `X-Request-ID` for trace correlation across services.
- Probes enable rolling updates/HPA in Kubernetes.
- Structured logs and metrics/tracing (OpenTelemetry-ready) can be added without changing handlers.

### 7) Delivery & Evolution
- Independent, small service with clear contracts → safer, faster CI/CD.
- Versioned APIs and message schemas support incremental evolution without breaking consumers.

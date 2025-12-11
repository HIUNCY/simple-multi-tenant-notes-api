# Multi-Tenant Notes API

A production-ready RESTful API service for managing notes with comprehensive multi-tenancy support, role-based access control, and asynchronous audit logging.

## Overview

This application provides a scalable notes management system built with Go, featuring organization-level data isolation, fine-grained permission control via Casbin, and distributed audit logging through message queuing architecture.

### Key Features

- **JWT Authentication**: Secure token-based authentication with user identity and organizational context
- **Multi-Tenancy**: Complete data isolation per organization using `org_id` segregation
- **Role-Based Access Control (RBAC)**: Policy-driven authorization powered by Casbin
- **Asynchronous Audit Logging**: Non-blocking event tracking via RabbitMQ message queue
- **Automated Database Migration**: Self-initializing schema on application startup
- **Production-Ready Containerization**: Multi-stage Docker builds with distroless base images

## Architecture

The application follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Gin HTTP  │────▶│  PostgreSQL  │     │  RabbitMQ   │
│   Server    │     │   (Notes)    │     │   (Queue)   │
└─────────────┘     └──────────────┘     └──────┬──────┘
      │                                         │
      │             ┌──────────────┐            │
      └────────────▶│   Casbin     │            │
                    │   (RBAC)     │            │
                    └──────────────┘            ▼
                                         ┌─────────────┐
                                         │   MongoDB   │
                                         │ (Audit Log) │
                                         └─────────────┘
```

### Technology Stack

- **Framework**: Gin (Go web framework)
- **Authorization**: Casbin (RBAC enforcement)
- **Primary Database**: PostgreSQL (relational data storage)
- **Audit Database**: MongoDB (document-based audit logs)
- **Message Broker**: RabbitMQ (AMQP protocol)
- **Containerization**: Docker (multi-stage builds with distroless runtime)

## Prerequisites

### Runtime Requirements

- Go 1.21+ (for local development)
- PostgreSQL 14+
- MongoDB 5.0+
- RabbitMQ 3.12+
- Docker 20.10+ (for containerized deployment)

### Development Tools

- Git
- Make (optional, for build automation)
- Docker Compose (optional, for local infrastructure)

## Configuration

The application uses environment variables for configuration. Create a `.env` file based on `.env.example`:

### Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_USER` | PostgreSQL username | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `your_password` |
| `DB_NAME` | PostgreSQL database name | `notes_db` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` |
| `MONGO_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `MONGO_DB_NAME` | MongoDB database name | `audit_logs` |
| `JWT_SECRET` | Secret key for JWT signing | `your_secret_key` |
| `RABBITMQ_URL` | RabbitMQ connection URL | `amqp://guest:guest@localhost:5672/` |

### Optional Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_CHANNEL_BINDING` | PostgreSQL channel binding | - |

## Installation & Deployment

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/HIUNCY/simple-multi-tenant-notes-api.git
   cd simple-multi-tenant-notes-api
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start dependencies**
   ```bash
   # Start PostgreSQL, MongoDB, and RabbitMQ
   # Using Docker Compose (recommended):
   docker-compose up -d postgres mongodb rabbitmq
   ```

4. **Run the application**
   ```bash
   go run ./cmd/api
   ```

The application will automatically:
- Connect to all configured services
- Run database migrations
- Start the RabbitMQ consumer for audit logging
- Begin accepting HTTP requests

### Docker Deployment

#### Build the Image

```bash
docker build -t notes-api:latest .
```

#### Run the Container

**Linux/macOS:**
```bash
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_USER=postgres \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=notes_db \
  -e DB_SSLMODE=disable \
  -e DB_PORT=5432 \
  -e MONGO_URI=mongodb://host.docker.internal:27017 \
  -e MONGO_DB_NAME=audit_logs \
  -e JWT_SECRET=your_secret_key \
  -e RABBITMQ_URL=amqp://guest:guest@host.docker.internal:5672/ \
  notes-api:latest
```

**Windows (PowerShell):**
```powershell
docker run --rm -p 8080:8080 `
  -e PORT=8080 `
  -e DB_HOST=host.docker.internal `
  -e DB_USER=postgres `
  -e DB_PASSWORD=your_password `
  -e DB_NAME=notes_db `
  -e DB_SSLMODE=disable `
  -e DB_PORT=5432 `
  -e MONGO_URI=mongodb://host.docker.internal:27017 `
  -e MONGO_DB_NAME=audit_logs `
  -e JWT_SECRET=your_secret_key `
  -e RABBITMQ_URL=amqp://guest:guest@host.docker.internal:5672/ `
  notes-api:latest
```

> **Note**: Replace `host.docker.internal` with appropriate service names when using Docker Compose or Kubernetes.

## API Documentation

### Authentication

All protected endpoints require a valid JWT token in the Authorization header.

#### Login

Obtain a JWT token for API access.

**Endpoint:** `POST /login`

**Request Body:**
```json
{
  "user_id": "user123",
  "org_id": "org456",
  "role": "admin"
}
```

**Response (200 OK):**
```json
{
  "message": "Login berhasil! Gunakan token ini di Header 'Authorization: Bearer <token>'"
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
}
```

**Usage:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user123","org_id":"org456","role":"admin"}'
```

### Notes Management

All notes endpoints require authentication via Bearer token.

#### Create Note

**Endpoint:** `POST /api/notes` (Admin only)

**Headers:**
```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Meeting Notes",
  "content": "Discussion points from today's meeting..."
}
```

**Response (201 Created):**
```json
{
  "data": {
    "id": "1",
    "title": "Meeting Notes",
    "content": "Discussion points from today's meeting...",
    "org_id": "org456",
    "user_id": "user123",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

#### List Notes

**Endpoint:** `GET /api/notes`

**Headers:**
```
Authorization: Bearer <your_jwt_token>
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid-string",
      "title": "Meeting Notes",
      "content": "Discussion points...",
      "org_id": "org456",
      "user_id": "user123",
      "created_at": "2025-01-15T10:30:00Z"
    }
  ]
}
```

#### Get Note by ID

**Endpoint:** `GET /api/notes/:id`

**Headers:**
```
Authorization: Bearer <your_jwt_token>
```

**Response (200 OK):**
```json
{
  "data": {
    "id": "uuid-string",
    "title": "Meeting Notes",
    "content": "Discussion points...",
    "org_id": "org456",
    "user_id": "user123",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

## Authorization & RBAC

The application uses Casbin for policy-based access control. Authorization rules are defined in:

- `model.conf`: RBAC model definition
- `policy.csv`: Role-permission mappings

### How It Works

The Casbin middleware evaluates each request based on:
- **Subject**: User role (from JWT token)
- **Object**: API endpoint path
- **Action**: HTTP method

Ensure that roles assigned during login match the policies defined in `policy.csv` to grant appropriate access.

### Example Policy

```csv
p, admin, /api/notes, POST
p, admin, /api/notes, GET
p, admin, /api/notes/:id, GET
p, user, /api/notes, GET
p, user, /api/notes/:id, GET
```

## Audit Logging

The application implements asynchronous audit logging for compliance and monitoring:

1. Business operations publish audit events to RabbitMQ
2. A dedicated consumer service reads from the queue
3. Events are persisted to MongoDB for long-term storage and analysis

This architecture ensures that audit logging does not impact API response times.

## Troubleshooting

### Common Issues

#### 401 Unauthorized

**Cause**: Missing or invalid JWT token

**Solution**: Ensure the `Authorization` header contains a valid Bearer token obtained from the `/login` endpoint. Verify that the token includes required claims (`user_id`, `org_id`, `role`).

#### 403 Forbidden

**Cause**: Insufficient permissions for the requested operation

**Solution**: Check that the user's role in the JWT token has appropriate permissions in `policy.csv`. Verify the Casbin policy configuration allows the role to access the specific endpoint and HTTP method.

#### Connection Failures

**Cause**: Unable to connect to PostgreSQL, MongoDB, or RabbitMQ

**Solution**: 
- Verify all database and message broker services are running
- Check environment variables for correct hostnames, ports, and credentials
- Ensure network connectivity between the application and services
- Review firewall rules and security groups

#### Port Already in Use

**Cause**: Another process is using the configured port

**Solution**: Change the `PORT` environment variable to an available port, or stop the conflicting process.

## Project Structure

```
.
├── cmd/
│   └── api/                    # Application entry point
├── internal/
│   ├── middleware/             # JWT and Casbin middleware
│   ├── repository/             # Data access layer
│   ├── service/                # Business logic
│   ├── queue/                  # RabbitMQ producer/consumer
│   └── utils/                  # JWT utilities
├── model.conf                  # Casbin RBAC model
├── policy.csv                  # Casbin authorization policies
├── Dockerfile                  # Multi-stage production build
├── .env.example                # Environment configuration template
└── README.md                   # Project documentation
```

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is intended for educational and demonstration purposes. Use and modify as needed for your requirements.

## Support

For issues, questions, or contributions, please open an issue on the project repository.

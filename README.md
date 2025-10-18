# Notifications Service

This service provides notification functionalities for ProofLift, consuming events from the posts service to deliver notifications to users.

## Running the Service

You can run this service either locally or using Docker.

### Option 1: Running Locally

#### 1. Clone the repository
```bash
git clone https://github.com/Swarch-2025ii-ProofLift/prooflift-notifications-be
cd prooflift-notifications-be
```

#### 2. Install Go dependencies
```bash
go mod download
```

#### 3. Set up environment variables
Create a `.env` file in the root of the project using `.env.example` as a template.

When running locally, set:

```
DB_HOST=localhost
MQ_HOST=localhost
```

This ensures the service connects to your locally installed PostgreSQL instance and RabbitMQ server.

#### 4. Start local dependencies
Ensure PostgreSQL and RabbitMQ are running locally:
- PostgreSQL should be accessible on port `5432`
- RabbitMQ should be accessible on port `5672`

#### 5. Run the service locally
```bash
go run cmd/app/main.go
```

The service will automatically run database migrations on startup.

### Option 2: Running with Docker

#### 1. Clone the repository
```bash
git clone https://github.com/Swarch-2025ii-ProofLift/prooflift-notifications-be
cd prooflift-notifications-be
```

#### 2. Set up environment variables
Create a `.env` file in the root of the project using `.env.example` as a template.

When running the service with Docker, the database and message queue hosts should point to the services defined in `docker-compose.yml`:

```
DB_HOST=prooflift-notifications-db
MQ_HOST=prooflift-notifications-mq
```

This ensures the service connects to the PostgreSQL and RabbitMQ containers inside the Docker network.

#### 3. Build and run with Docker Compose
```bash
docker-compose up --build
```

### Accessing the API

The API will be available at (both locally and with Docker):
- API Service: `http://localhost:8080`
- Health check: `http://localhost:8080/health`
- RabbitMQ Management UI (Docker only): `http://localhost:15672` (default credentials: guest/guest)

## Message Queue Integration

This service consumes notification events from a RabbitMQ message queue. The posts service publish events when certain actions occur:
- **Comment events**: When someone comments on a post
- **Reaction events**: When someone reacts to a post

The service processes these events and creates notifications in the database for the relevant users.

The MQ connection is configured via environment variables (`MQ_HOST`, `MQ_PORT`, `MQ_USER`, `MQ_PASSWORD`, `QUEUE_NAME`).

## Project Structure
```
cmd/
└── app/
    └── main.go                 # Entry point

internal/
├── api/
│   ├── handlers/               # HTTP request handlers
│   ├── http/                   # HTTP router setup
│   ├── middleware/             # Middleware elements
│   ├── responses/              # Response utilities
│   └── schemas/                # Request/response schemas
│
├── configs/                    # Configuration management
│
├── db/                         # Database connection and migrations
│
├── events/                     # Event consumer (MQ)
│
├── models/                     # Database models
│
├── mq/                         # Message queue connection
│
├── repositories/               # Data access layer
│
└── services/                   # Business logic
```

## API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/health` | Service health check | No |
| `GET` | `/api/notifications` | Get user's notifications | Yes (JWT) |
| `PUT` | `/api/notifications/{id}/read` | Mark a notification as read | Yes (JWT) |
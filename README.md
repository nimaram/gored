# Gored

A Go-based microservice application that exposes an HTTP API, publishes messages to Redpanda (Kafka-compatible message broker), and processes them with a background worker. Traffic is routed through an Nginx reverse proxy and protected by a simple in-memory rate limiter.

## Project Structure

```text
gored/
├── main.go                 # HTTP server entry point (Gin)
├── worker/
│   └── main.go             # Worker service entry point
├── services/
│   ├── producer.go         # Kafka/Redpanda producer service
│   └── worker.go           # Kafka/Redpanda consumer/worker
├── utils/
│   └── ratelimit/
│       └── leaky_bucket.go # Gin rate limiting middleware
├── nginx.conf              # Nginx reverse proxy configuration
├── compose.yaml            # Docker Compose (app + worker + Redpanda + Nginx)
├── Dockerfile              # Multi-stage Docker build (app + worker binaries)
├── .env-sample             # Sample env file for Docker Compose
├── go.mod                  # Go module dependencies
└── README.md               # This file
```

## Architecture

The project follows a simple event-driven microservice architecture:

```text
   Client
     │
     ▼
┌────────────┐        ┌─────────────────┐
│   Nginx    │ 80     │   HTTP Server   │ 8080
│ (nginx.conf) ─────▶ │   (main.go)     │
└────────────┘        │  + Rate limiting│
                       └────────┬────────┘
                                │
                                │ Publishes
                                ▼
                       ┌─────────────────┐
                       │  Producer       │ (services/producer.go)
                       └────────┬────────┘
                                │
                                │ Writes to
                                ▼
                       ┌─────────────────┐
                       │   Redpanda      │ (Kafka-compatible broker)
                       │   (compose)     │
                       └────────┬────────┘
                                │
                                │ Consumes
                                ▼
                       ┌─────────────────┐
                       │   Worker        │ (services/worker.go + worker/main.go)
                       └─────────────────┘
```

### Components

- **HTTP API Server** (`main.go`):
  - Built with Gin web framework.
  - Uses a leaky-bucket rate limiting middleware (`utils/ratelimit/leaky_bucket.go`).
  - Exposes:
    - `GET /ping` – basic JSON response.
    - `GET /healthz` – health check endpoint (used by Docker healthcheck).
    - `POST /task` – enqueues a task message into Redpanda.
  - Listens on port `8080` inside the container (fronted by Nginx on port `80`).

- **Producer Service** (`services/producer.go`):
  - Handles message publishing to Redpanda using `segmentio/kafka-go`.
  - Publishes messages to the `tasks` topic.

- **Worker Service** (`services/worker.go`, `worker/main.go`):
  - Runs as a separate binary/container.
  - Consumes messages from the `tasks` topic.
  - Simulates processing work with a short delay and logs progress.

- **Nginx Reverse Proxy** (`nginx.conf`):
  - Listens on port `80`.
  - Proxies incoming HTTP traffic to the app service on `app:8080`.

- **Redpanda** (via `compose.yaml`):
  - Kafka-compatible message broker.
  - Exposed on port `9092` (inside and outside Docker).
  - Configured for lightweight local development.

## Tools & Technologies

- **Go 1.24.2**: Programming language.
- **Gin**: HTTP web framework for building REST APIs.
- **Redpanda**: Kafka-compatible streaming data platform.
- **segmentio/kafka-go**: Go library for Kafka/Redpanda integration.
- **Nginx**: Reverse proxy in front of the HTTP API.
- **Docker & Docker Compose**: Containerization and orchestration.

## Prerequisites

- Go 1.24.2 or later (for local development).
- Docker and Docker Compose.
- Available ports:
  - `80` (Nginx reverse proxy)
  - `9092` (Redpanda)

## Configuration

### Environment Variables

- **`REDPANDA_STRING_URL`**: Connection string for Redpanda broker.
  - Inside Docker network, typically: `redpanda:9092`.
  - From host/local development, typically: `localhost:9092`.

### `.env` file (Docker Compose)

For Docker Compose, you can configure environment variables using a `.env` file in the project root. A sample is provided:

```bash
cp .env-sample .env
```

Edit `.env` and set:

```bash
REDPANDA_STRING_URL=redpanda:9092
```

This value is injected into both the `app` and `worker` services.

## How to Run

### Option 1: Using Docker Compose (Recommended)

This option runs the full stack: Redpanda, HTTP app, worker, and Nginx.

1. **Create `.env` from sample (once):**

   ```bash
   cp .env-sample .env
   # then edit .env if needed
   ```

2. **Start all services:**

   ```bash
   docker compose up --build
   ```

   This will start:
   - `redpanda` (message broker)
   - `app` (HTTP API service)
   - `worker` (background worker)
   - `nginx` (reverse proxy on port 80)

3. **Test the API (through Nginx on port 80):**

   ```bash
   curl http://localhost/ping
   curl -X POST http://localhost/task
   ```

4. **Check health endpoint (used by Compose healthcheck):**

   ```bash
   curl http://localhost/healthz
   ```

### Option 2: Local Development (without Docker for app/worker)

In this mode, Redpanda still runs via Docker, but the Go services run directly on your machine.

1. **Start Redpanda only:**

   ```bash
   docker compose up -d redpanda
   ```

2. **Install Go dependencies:**

   ```bash
   go mod download
   ```

3. **Set environment variable for Redpanda:**

   ```bash
   # Windows (PowerShell)
   $env:REDPANDA_STRING_URL = "localhost:9092"

   # Linux / macOS (bash/zsh)
   export REDPANDA_STRING_URL="localhost:9092"
   ```

4. **Run the HTTP app:**

   ```bash
   go run main.go
   ```

5. **Run the worker in a separate terminal:**

   ```bash
   go run ./worker
   ```

6. **Test the API directly (bypassing Nginx):**

   ```bash
   curl http://localhost:8080/ping
   curl -X POST http://localhost:8080/task
   ```

### Option 3: Build Image Manually (advanced)

You can also build the image manually if you want to run containers without Compose coordination:

1. **Build image:**

   ```bash
   docker build -t gored .
   ```

2. **Ensure Redpanda is running (e.g., via Compose):**

   ```bash
   docker compose up -d redpanda
   ```

3. **Run the HTTP app container:**

   ```bash
   docker run -p 8080:8080 -e REDPANDA_STRING_URL=host.docker.internal:9092 gored
   ```

   On Linux you may need to replace `host.docker.internal` with the Docker bridge IP (e.g. `172.17.0.1`) or your host IP.

> Note: With this option you are responsible for running the worker and reverse proxy yourself (Compose is the easiest way to run the full stack).

## API Endpoints

All endpoints are available:
- via Nginx at `http://localhost` when using full Docker Compose.
- directly at `http://localhost:8080` when running the Go app locally without Nginx.

- **`GET /ping`**
  - Response example:
    ```json
    {"message": "Hello World!"}
    ```

- **`GET /healthz`**
  - Returns HTTP 200 when the app is healthy.

- **`POST /task`**
  - Enqueues a task into Redpanda on the `tasks` topic.
  - Example:
    ```bash
    curl -X POST http://localhost/task
    ```
  - Returns:
    ```json
    {"status": "queued"}
    ```

### Rate Limiting

- A simple in-memory leaky bucket rate limiter is applied per-client-IP.
- Default configuration:
  - Capacity: 10 tokens.
  - Refill: 1 token per second.
- Excess requests are rejected with HTTP `429 Too Many Requests`.

## Stopping Services

To stop all Docker services started via Compose:

```bash
docker compose down
```

## Development Notes

- This project is intended as a simple example of:
  - HTTP API with rate limiting.
  - Event-driven processing using Redpanda (Kafka-compatible).
  - Separating HTTP and worker responsibilities into different processes.
- The configuration is optimized for local development and experimentation, not production hardening.

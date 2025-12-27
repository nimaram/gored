# Gored

A Go-based microservice application that provides an HTTP API and publishes messages to Redpanda (Kafka-compatible message broker).

## Project Structure

```
gored/
├── main.go              # HTTP server entry point
├── services/
│   └── producer.go      # Kafka/Redpanda producer service
├── compose.yaml         # Docker Compose configuration for Redpanda
├── Dockerfile           # Multi-stage Docker build
├── go.mod              # Go module dependencies
└── README.md           # This file
```

## Architecture

The project follows a simple microservice architecture:

```
┌─────────────────┐
│   HTTP Server   │  (Gin Framework)
│   (main.go)     │
└────────┬────────┘
         │
         │ Uses
         ▼
┌─────────────────┐
│  Producer       │  (services/producer.go)
│  Service        │
└────────┬────────┘
         │
         │ Publishes to
         ▼
┌─────────────────┐
│   Redpanda      │  (Kafka-compatible broker)
│   (Docker)      │
└─────────────────┘
```

### Components

- **HTTP API Server** (`main.go`): 
  - Built with Gin web framework
  - Exposes REST endpoints (currently `/ping`)
  - Runs on port 8080 (default Gin port)

- **Producer Service** (`services/producer.go`):
  - Handles message publishing to Redpanda
  - Uses `segmentio/kafka-go` library
  - Publishes messages to the `tasks` topic
  - Requires `REDPANDA_STRING_URL` environment variable

- **Redpanda** (via Docker Compose):
  - Kafka-compatible message broker
  - Exposed on port 9092
  - Configured for development with minimal resources

## Tools & Technologies

- **Go 1.24.2**: Programming language
- **Gin**: HTTP web framework for building REST APIs
- **Redpanda**: Kafka-compatible streaming data platform
- **segmentio/kafka-go**: Go library for Kafka/Redpanda integration

## Prerequisites

- Go 1.24.2 or later
- Docker and Docker Compose
- Make sure ports 8080 and 9092 are available

## How to Run

### Option 1: Using Docker Compose (Recommended)

1. **Start Redpanda service:**
   ```bash
   docker-compose up -d
   ```

2. **Set environment variable:**
   ```bash
   # On Windows (PowerShell)
   $env:REDPANDA_STRING_URL="localhost:9092"
   
   # On Linux/Mac
   export REDPANDA_STRING_URL="localhost:9092"
   ```

3. **Run the application:**
   ```bash
   go run main.go
   ```

4. **Test the API:**
   ```bash
   curl http://localhost:8080/ping
   ```

### Option 2: Local Development

1. **Start Redpanda:**
   ```bash
   docker-compose up -d
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set environment variable:**
   ```bash
   # On Windows (PowerShell)
   $env:REDPANDA_STRING_URL="localhost:9092"
   
   # On Linux/Mac
   export REDPANDA_STRING_URL="localhost:9092"
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```

### Option 3: Using Docker

1. **Build the Docker image:**
   ```bash
   docker build -t gored .
   ```

2. **Start Redpanda:**
   ```bash
   docker-compose up -d
   ```

3. **Run the container:**
   ```bash
   docker run -p 8080:8080 -e REDPANDA_STRING_URL=host.docker.internal:9092 gored
   ```

   Note: On Linux, you may need to use `172.17.0.1:9092` or the actual Redpanda container IP instead of `host.docker.internal`.

## Environment Variables

- `REDPANDA_STRING_URL`: Connection string for Redpanda broker (e.g., `localhost:9092` or `redpanda:9092` when running in Docker network)

## API Endpoints

### GET /ping
Returns a simple JSON response.


## Stopping Services

To stop Redpanda:
```bash
docker-compose down
```

## Development Notes

- It's currently under development.
- The producer service writes to the `tasks` topic in Redpanda
- The HTTP server runs on the default Gin port (8080)
- Redpanda is configured for development with minimal resources (1 CPU, 1GB memory)

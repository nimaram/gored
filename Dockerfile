# --------- build stage ----------
    FROM golang:1.24.2-alpine AS builder

    WORKDIR /app
    
    # Cache dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy source
    COPY . .
    
    # Build HTTP app binary from root main.go
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -o gored .
    
    # Build worker binary from /worker/main.go
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -o worker ./worker
    
    
    # ---------- runtime stage ----------
    FROM alpine:latest AS runtime

    WORKDIR /app

    RUN apk --no-cache add curl wget
    
    # Copy binaries from builder image
    COPY --from=builder /app/gored /app/gored
    COPY --from=builder /app/worker /app/worker
    
    # App listens on 8080
    EXPOSE 8080
    
    RUN adduser -D nonroot
    USER nonroot:nonroot
    
    # Default entrypoint is the HTTP app; docker-compose will override for worker
    ENTRYPOINT ["/app/gored"]
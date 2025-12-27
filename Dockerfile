# --------- build stage ----------
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./

RUN go mod download

# copy source
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o gored ./cmd/gored


# ---------- runtime stage ----------
FROM gcr.io/distroless/base-debian12 AS runtime

WORKDIR /app

COPY --from=builder /app/gored /app/gored

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/gored"]
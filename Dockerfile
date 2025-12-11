## Multi-stage build for Gin notes API
# Build stage
FROM golang:1.23-bullseye AS builder

WORKDIR /app

# Enable Go modules and better caching of deps
ENV GO111MODULE=on \
    CGO_ENABLED=0

# Pre-copy go.mod and go.sum to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the API binary
RUN go build -o /out/app ./cmd/api


# Runtime stage
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Copy executable
COPY --from=builder /out/app /app/app

# Copy Casbin model and policy files required at runtime
COPY model.conf policy.csv /app/

# Expose default port (override with PORT env as needed)
EXPOSE 8080

# Run as non-root user
USER nonroot:nonroot

# Start the server
ENTRYPOINT ["/app/app"]

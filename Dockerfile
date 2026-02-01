# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web

# Copy package files
COPY web/package*.json ./

# Install dependencies (use ci for reproducible builds)
RUN npm ci --only=production=false

# Copy source code
COPY web/ .

# Build frontend
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Install swag for Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/aigc-check-server/main.go -o docs

# Build backend
# CGO_ENABLED=1 is required for go-sqlite3
# Use musl-compatible build flags for Alpine
RUN CGO_ENABLED=1 GOOS=linux \
    CGO_CFLAGS="-D_LARGEFILE64_SOURCE" \
    go build \
    -ldflags="-w -s" \
    -o aigc-check-server cmd/aigc-check-server/main.go

# Stage 3: Final Image
FROM alpine:3.19
WORKDIR /app

# Install runtime dependencies and create non-root user
RUN apk add --no-cache sqlite-libs ca-certificates tzdata && \
    adduser -D -u 1000 appuser && \
    mkdir -p /app/data && \
    chown -R appuser:appuser /app

# Set timezone
ENV TZ=Asia/Shanghai

# Copy backend binary
COPY --from=backend-builder /app/aigc-check-server .

# Copy config files
COPY --from=backend-builder /app/configs ./configs

# Copy frontend build artifacts to static directory
COPY --from=frontend-builder /app/web/dist ./static

# Set ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the server
CMD ["./aigc-check-server"]

# Stage 1: Build the application
FROM golang:1.24.3-alpine3.21 AS builder

# Install build dependencies
RUN apk add --no-cache  \
  ca-certificates  \
  git \
  tzdata

WORKDIR /app

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./

# Download dependencies (this layer will be cached if go.mod/go.sum don't change)
RUN go mod download

# Install sqlc for code generation
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install stringer
RUN go install golang.org/x/tools/cmd/stringer@latest

# Copy source code
COPY ./main.go ./main.go
COPY ./config ./config
COPY ./internal ./internal
COPY ./queries ./queries
COPY ./migrations ./migrations
COPY ./sqlc.yaml ./sqlc.yaml

# Run go generate to generate any required code
RUN go generate ./...

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags='-w -s -extldflags "-static"' \
  -a -installsuffix cgo \
  -o server .

# Stage 2: Create runtime image with Chrome
# Option 1: Alpine with Chromium (smaller image)
FROM alpine:3.21

LABEL com.centurylinklabs.watchtower.enable="true"

WORKDIR /app

# Install Chrome and dependencies
RUN apk add --no-cache \
  ca-certificates \
  chromium \
  chromium-chromedriver \
  dumb-init \
  tzdata

# Set Chrome binary path environment variable
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser

# Copy timezone data and certificates
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the configuration file
COPY ./config ./config

# Copy the migration files
COPY ./migrations ./migrations

# Copy the binary
COPY --from=builder /app/server ./server


RUN mkdir -p ./logs && \
  addgroup -g 1001 -S appgroup && \
  adduser -u 1001 -S appuser -G appgroup && \
  chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

ARG CONFIG_FILE
ENV CONFIG_FILE ${CONFIG_FILE}

# Use dumb-init to handle signals properly
ENTRYPOINT ["dumb-init", "--"]

# Run the server
CMD ["./server"]

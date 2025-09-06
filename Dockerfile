# Multi-stage build for GVM SSH Setup Tool
# Optimized for private registry deployment

# Build stage
FROM golang:1.22-alpine AS build

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /src

# Copy dependency files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /out/gvm-ssh \
    ./main.go

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    openssh-client \
    git \
    bash \
    ca-certificates \
    curl

# Create non-root user
RUN adduser -D -s /bin/bash -u 1000 gvmuser

# Copy binary from build stage
COPY --from=build /out/gvm-ssh /usr/local/bin/gvm-ssh

# Copy entrypoint script
COPY entrypoint.sh /usr/local/bin/entrypoint

# Set permissions
RUN chmod +x /usr/local/bin/gvm-ssh /usr/local/bin/entrypoint

# Set up environment
ENV PATH="/usr/local/bin:${PATH}"

# Use entrypoint
ENTRYPOINT ["entrypoint"]

# Default to wizard mode
CMD ["wizard"]

# Labels for metadata
LABEL org.opencontainers.image.title="GVM SSH Setup Tool"
LABEL org.opencontainers.image.description="SSH and Git configuration tool for GitLab CE environments"
LABEL org.opencontainers.image.vendor="Global Vision Media"
LABEL org.opencontainers.image.source="https://github.com/hhawkinsgvm/gvm-ssh-setup"
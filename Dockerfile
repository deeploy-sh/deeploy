# Build-Stage
FROM golang:1.24-alpine3.20 AS build
WORKDIR /app

# Install build dependencies first (rarely changes)
RUN apk add --no-cache gcc musl-dev wget

# Install templ (rarely changes)
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy dependency files first (changes less often than code)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (changes most often - do this last)
COPY . .

# Generate templ files
RUN templ generate

# Install Tailwind CSS standalone CLI (musl for Alpine)
RUN ARCH=$(uname -m) && \
  if [ "$ARCH" = "x86_64" ]; then \
    TAILWIND_URL="https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.3/tailwindcss-linux-x64-musl"; \
  elif [ "$ARCH" = "aarch64" ]; then \
    TAILWIND_URL="https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.3/tailwindcss-linux-arm64-musl"; \
  else \
    echo "Unsupported architecture: $ARCH"; exit 1; \
  fi && \
  wget -O tailwindcss "$TAILWIND_URL" && chmod +x tailwindcss

# Generate Tailwind CSS (must happen before go build for embed)
RUN ./tailwindcss -i ./internal/server/assets/css/input.css -o ./internal/server/assets/css/output.css --minify

# Build Go binary
ARG VERSION=dev
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-X github.com/deeploy-sh/deeploy/internal/shared/version.Version=$VERSION" \
    -o main ./cmd/server

# Deploy-Stage
FROM alpine:3.21
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates git

# Set environment variable for runtime
ENV APP_ENV=production

# Copy the binary from the build stage
COPY --from=build /app/main .

# Expose the port your application runs on
EXPOSE 8090

# Command to run the application
CMD ["./main"]

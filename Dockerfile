FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Copy source code
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Cambiar la ruta del build para apuntar a cmd/
RUN go build -o main ./cmd/server/main.go

# Final image with Chrome
FROM alpine:latest

# Install Chrome and dependencies
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

# Configure Chrome for container execution
ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/bin/chromium-browser

# Create non-root user for security
RUN addgroup -g 1001 -S golang && \
    adduser -S golang -u 1001

# Copy binary
COPY --from=builder /app/main /app/main
RUN chown golang:golang /app/main

USER golang
EXPOSE 3333

CMD ["/app/main"]


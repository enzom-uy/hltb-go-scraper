FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Cambiar 'main' por 'server' para coincidir con Railway
RUN go build -o server ./cmd/server

FROM alpine:latest

RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/bin/chromium-browser

RUN addgroup -g 1001 -S golang && \
    adduser -S golang -u 1001

# Cambiar 'main' por 'server'
COPY --from=builder /app/server /app/server
RUN chown golang:golang /app/server

USER golang
EXPOSE 3333

# Cambiar el comando para ejecutar 'server'
CMD ["/app/server"]


FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

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

COPY --from=builder /app/main /app/main
RUN chown golang:golang /app/main

USER golang
EXPOSE 3333

CMD ["/app/main"]

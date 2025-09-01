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

USER golang
EXPOSE 3333

# Railway will handle the build and start commands
# No CMD needed - Railway uses your custom start command

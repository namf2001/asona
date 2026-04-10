ARG GO_VERSION=1.23-alpine

# ── Dev container with Air hot-reload ─────────────────────────────────────────
FROM golang:${GO_VERSION} AS dev

WORKDIR /app

# Install Air for hot-reload
RUN go install github.com/air-verse/air@latest

# Pre-download modules (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Source is mounted via volume at runtime, not copied
CMD ["air", "-c", ".air.toml"]

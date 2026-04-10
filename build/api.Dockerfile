ARG GO_VERSION=1.23-alpine
ARG RELEASE_IMAGE_NAME=alpine
ARG RELEASE_IMAGE_TAG=3.20

# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:${GO_VERSION} AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main cmd/api/main.go

# ── Stage 2: Production ───────────────────────────────────────────────────────
FROM ${RELEASE_IMAGE_NAME}:${RELEASE_IMAGE_TAG} AS prod

WORKDIR /app

RUN apk --no-cache add tzdata ca-certificates

COPY --from=builder /app/main /app/main

EXPOSE ${PORT}

CMD ["./main"]

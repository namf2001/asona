FROM golang:1.24.4-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


FROM node:20-alpine AS frontend_builder
WORKDIR /frontend

RUN corepack enable && corepack prepare pnpm@latest --activate

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY frontend/. .
RUN pnpm run build

FROM node:20-alpine AS frontend
RUN corepack enable && corepack prepare pnpm@latest --activate && pnpm add -g serve
COPY --from=frontend_builder /frontend/.next /app/.next
COPY --from=frontend_builder /frontend/public /app/public
COPY --from=frontend_builder /frontend/package.json /app/package.json
WORKDIR /app
EXPOSE 3000
CMD ["serve", "-s", "/app/.next", "-l", "3000"]

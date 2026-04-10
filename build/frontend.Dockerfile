ARG NODE_VERSION=20-alpine

# ── Stage 1: Install dependencies ─────────────────────────────────────────────
FROM node:${NODE_VERSION} AS deps

WORKDIR /frontend

RUN corepack enable && corepack prepare pnpm@latest --activate

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# ── Stage 2: Build ────────────────────────────────────────────────────────────
FROM node:${NODE_VERSION} AS builder

WORKDIR /frontend

RUN corepack enable && corepack prepare pnpm@latest --activate

COPY --from=deps /frontend/node_modules ./node_modules
COPY . .

RUN pnpm run build

# ── Stage 3: Production ───────────────────────────────────────────────────────
FROM node:${NODE_VERSION} AS prod

WORKDIR /app

# Copy standalone output from Next.js (output: 'standalone' in next.config.ts)
COPY --from=builder /frontend/.next/standalone ./
COPY --from=builder /frontend/.next/static ./.next/static
COPY --from=builder /frontend/public ./public

EXPOSE 3000

ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]

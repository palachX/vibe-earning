# ─── СТАДИЯ 1: СБОРКА ─────────────────────────────────────────────────────────
FROM bitnami/node:latest AS builder

WORKDIR /app
RUN corepack enable && corepack prepare pnpm@latest --activate

COPY package.json pnpm-lock.yaml ./
COPY .env ./

RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm build

# ─── СТАДИЯ 2: РАНТАЙМ ────────────────────────────────────────────────────────
FROM bitnami/node:latest AS runtime

WORKDIR /app

COPY --from=builder /app/.output/ ./

ENV PORT 3000
ENV HOST 0.0.0.0

EXPOSE 3000
ENTRYPOINT ["node", "/app/server/index.mjs", "--host", "0.0.0.0", "--port", "3000"]


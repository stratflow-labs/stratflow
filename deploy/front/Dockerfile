# syntax=docker/dockerfile:1

FROM node:24-alpine AS base

WORKDIR /app

ENV NEXT_TELEMETRY_DISABLED=1

COPY package.json package-lock.json ./
COPY apps/web/package.json ./apps/web/package.json

RUN npm install --no-audit --no-fund


FROM base AS dev

WORKDIR /app

ENV APP_ENV=dev
ENV NODE_ENV=development

COPY . .

WORKDIR /app/apps/web

EXPOSE 3000

CMD ["npx", "next", "dev", "--hostname", "0.0.0.0", "--port", "3000"]


FROM base AS builder

WORKDIR /app

ENV NODE_ENV=production

ARG NEXT_PUBLIC_API_BASE_URL
ENV NEXT_PUBLIC_API_BASE_URL=${NEXT_PUBLIC_API_BASE_URL}

COPY . .

RUN npm run build -w @stratflow/web


FROM node:24-alpine AS prod

WORKDIR /app

ENV APP_ENV=prod
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

COPY --from=builder /app/package.json /app/package-lock.json ./
COPY --from=builder /app/apps/web/package.json ./apps/web/package.json
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/apps/web/.next ./apps/web/.next
COPY --from=builder /app/apps/web/public ./apps/web/public
COPY --from=builder /app/apps/web/next.config.* ./apps/web/

WORKDIR /app/apps/web

EXPOSE 3000

CMD ["npx", "next", "start", "--hostname", "0.0.0.0", "--port", "3000"]
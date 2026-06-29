# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS base

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download


FROM base AS dev

WORKDIR /app

ENV APP_ENV=dev

RUN go install github.com/air-verse/air@latest

COPY . .

CMD ["air"]


FROM base AS builder

ARG SERVICE

COPY . .

RUN test -n "$SERVICE"

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/app \
    ./services/${SERVICE}/cmd/api


FROM alpine:3.22 AS prod

WORKDIR /app

RUN apk add --no-cache ca-certificates \
    && adduser -D -H -u 10001 appuser

COPY --from=builder /out/app /app/app

USER appuser

ENTRYPOINT ["/app/app"]
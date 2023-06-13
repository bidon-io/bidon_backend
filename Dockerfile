# Build UI
FROM node:18-alpine AS frontend-deps

WORKDIR /app

COPY web/bidon_ui/package.json web/bidon_ui/yarn.lock ./
RUN yarn install

FROM frontend-deps AS frontend-builder

ARG APP_ENV=production
COPY web/bidon_ui .
RUN VITE_APP_ENV=${APP_ENV} yarn generate

FROM golang:1.20-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal

FROM base AS test

CMD [ "go", "test", "./..." ]

FROM base AS builder

ARG BIDON_SERVICE

COPY --from=frontend-builder /app/.output/public ./cmd/$BIDON_SERVICE/web/ui
RUN go build -o /$BIDON_SERVICE ./cmd/$BIDON_SERVICE

FROM alpine:3.18

ARG BIDON_SERVICE

RUN adduser -D -u 1000 deploy
USER deploy

COPY --from=builder --chown=deploy /$BIDON_SERVICE /bidon-service

EXPOSE 1323

CMD [ "/bidon-service" ]

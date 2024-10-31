# Build UI
FROM node:22-alpine AS frontend-deps

WORKDIR /app

COPY web/bidon_ui/package.json web/bidon_ui/yarn.lock ./
RUN yarn install --frozen-lockfile

ARG APP_ENV=production
COPY web/bidon_ui .
RUN VITE_APP_ENV=${APP_ENV} yarn generate

FROM golang:1.23-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM base AS pre-commit-deps

RUN apk add --no-cache python3 git pre-commit

FROM base AS bidon-admin-builder

COPY --from=frontend-deps /app/.output/public ./cmd/bidon-admin/web/ui
RUN go build -o /bidon-admin ./cmd/bidon-admin

FROM base AS bidon-sdkapi-builder

RUN go build -o /bidon-sdkapi ./cmd/bidon-sdkapi

FROM base AS bidon-migrate-builder

RUN go build -o /bidon-migrate ./cmd/bidon-migrate

FROM alpine:3.18 AS deploy

RUN apk add --no-cache ca-certificates

RUN adduser -D -u 1000 deploy
USER deploy

EXPOSE 1323

FROM deploy AS bidon-admin

COPY --from=bidon-admin-builder --chown=deploy /bidon-admin /bidon-admin

CMD [ "/bidon-admin" ]

FROM deploy AS bidon-sdkapi

COPY --from=bidon-sdkapi-builder --chown=deploy /bidon-sdkapi /bidon-sdkapi

CMD [ "/bidon-sdkapi" ]

FROM deploy AS bidon-migrate

COPY --from=bidon-migrate-builder --chown=deploy /bidon-migrate /bidon-migrate

ENTRYPOINT [ "/bidon-migrate" ]

CMD [ "status" ]

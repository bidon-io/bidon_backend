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

# Run tests for each package sequantially, because each test package that accesses the database runs database.AutoMigrate.
# Fix this by running migrations before tests as a separate step.
CMD [ "go", "test", "-p", "1", "./..." ]

FROM base AS bidon-admin-builder

COPY --from=frontend-builder /app/.output/public ./cmd/bidon-admin/web/ui
RUN go build -o /bidon-admin ./cmd/bidon-admin

FROM base AS bidon-sdkapi-builder

RUN go build -o /bidon-sdkapi ./cmd/bidon-sdkapi

FROM alpine:3.18 AS deploy

RUN adduser -D -u 1000 deploy
USER deploy

EXPOSE 1323

FROM deploy AS bidon-admin

COPY --from=bidon-admin-builder --chown=deploy /bidon-admin /bidon-admin

CMD [ "/bidon-admin" ]

FROM deploy AS bidon-sdkapi

COPY --from=bidon-sdkapi-builder --chown=deploy /bidon-sdkapi /bidon-sdkapi

CMD [ "/bidon-sdkapi" ]

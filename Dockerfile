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

RUN go build -o /$BIDON_SERVICE ./cmd/$BIDON_SERVICE

FROM alpine:3.18

ARG BIDON_SERVICE

RUN adduser -D -u 1000 deploy
USER deploy

COPY --from=builder --chown=deploy /$BIDON_SERVICE /bidon-service

EXPOSE 1323

CMD [ "/bidon-service" ]

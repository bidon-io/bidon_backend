FROM golang:alpine AS base

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

FROM alpine

ARG BIDON_SERVICE

RUN adduser -D -u 1000 app
USER app

COPY --from=builder --chown=app /$BIDON_SERVICE /bidon-service

EXPOSE 1323

CMD [ "/bidon-service" ]

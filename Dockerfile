# syntax=docker/dockerfile:1

FROM golang:1.23.8-alpine AS builder

ARG APP_VERSION="undefined"
ARG BUILD_TIME="undefined"

WORKDIR /go/src/github.com/artarts36/sentry-notifier

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -ldflags="-s -w" -o /go/bin/sentry-notifier /go/src/github.com/artarts36/sentry-notifier/cmd/main.go

######################################################

FROM alpine

COPY --from=builder /go/bin/sentry-notifier /go/bin/sentry-notifier

# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.title="sentry-notifier"
LABEL org.opencontainers.image.description="sentry-notifier"
LABEL org.opencontainers.image.url="https://github.com/artarts36/sentry-notifier"
LABEL org.opencontainers.image.source="https://github.com/artarts36/sentry-notifier"
LABEL org.opencontainers.image.vendor="ArtARTs36"
LABEL org.opencontainers.image.version="$APP_VERSION"
LABEL org.opencontainers.image.created="$BUILD_TIME"
LABEL org.opencontainers.image.licenses="MIT"

EXPOSE 8080

WORKDIR /app

CMD ["/go/bin/sentry-notifier"]

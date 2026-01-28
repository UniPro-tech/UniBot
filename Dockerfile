FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk update && apk add --no-cache \
    curl git file unzip \
    python3 make g++ ffmpeg \
    pkgconf \
    opus-dev \
    opusfile-dev

COPY . .

WORKDIR /app/src
RUN go mod tidy

RUN go build -o /app/src/main ./cmd/bot

FROM alpine:latest
WORKDIR /root/

RUN apk update && apk add --no-cache \
    opus \
    opus-dev \
    opusfile \
    opusfile-dev \
    ffmpeg

COPY --from=builder /app/src/main .

CMD ["./main"]

# OCI Metadata

LABEL org.opencontainers.image.authors="Yuito Akatsuki <yuito@yuito-it.jp>"
LABEL org.opencontainers.image.url="https://github.com/UniPro-tech/UniBot"
LABEL org.opencontainers.image.source="https://github.com/UniPro-tech/UniBot"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="UniBot"
LABEL org.opencontainers.image.description="A multifunctional Discord bot for community management, entertainment, and productivity."
LABEL org.opencontainers.image.vendor="All-Japan Digital Creative Club UniProject <info@uniproject.jp>"

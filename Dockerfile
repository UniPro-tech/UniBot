# check=error=true
# syntax=docker/dockerfile:1
FROM golang:1.26.2-alpine3.22 AS builder
WORKDIR /app

RUN apk update && apk add --no-cache \
    curl git file unzip \
    python3 make g++ ffmpeg \
    pkgconf \
    opus-dev \
    opusfile-dev

WORKDIR /app/src

RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x
COPY . .

RUN ../scripts/_build.sh

FROM gcr.io/distroless/static-debian13:nonroot
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

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

RUN git clone https://github.com/disgoorg/godave.git

ENV SHELL=/bin/sh
ENV VCPKG_FORCE_SYSTEM_BINARIES=1
ENV CC=/usr/bin/gcc CXX=/usr/bin/g++
ENV CXXFLAGS="-Wno-error=maybe-uninitialized"
RUN apk add build-base cmake ninja zip unzip curl git pkgconfig perl nasm go
RUN FORCE_BUILD=1 ./godave/scripts/libdave_install.sh v1.1.0
ENV PKG_CONFIG_PATH=/root/.local/lib/pkgconfig
RUN rm -r ./godave

COPY . .

WORKDIR /app/src

RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
    go mod download -x

RUN ../scripts/_build.sh

FROM alpine:3.23.4
WORKDIR /root/

RUN apk update && apk add --no-cache \
    opus \
    opus-dev \
    opusfile \
    opusfile-dev \
    ffmpeg

ENV PKG_CONFIG_PATH=/root/.local/lib/pkgconfig
ENV LD_LIBRARY_PATH=/root/.local/lib
COPY --from=builder /root/.local/ /root/.local/

COPY --from=builder /app/src/main .
RUN chmod 777 ./main

# non root
USER nobody

CMD ["./main"]

# OCI Metadata

LABEL org.opencontainers.image.authors="Yuito Akatsuki <yuito@yuito-it.jp>"
LABEL org.opencontainers.image.url="https://github.com/UniPro-tech/UniBot"
LABEL org.opencontainers.image.source="https://github.com/UniPro-tech/UniBot"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="UniBot"
LABEL org.opencontainers.image.description="A multifunctional Discord bot for community management, entertainment, and productivity."
LABEL org.opencontainers.image.vendor="All-Japan Digital Creative Club UniProject <info@uniproject.jp>"

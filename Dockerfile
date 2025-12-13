FROM node:20-bullseye

LABEL org.opencontainers.image.authors="Yuito Akatsuki <yuito@yuito-it.jp>"
LABEL org.opencontainers.image.url="https://github.com/UniPro-tech/UniBot"
LABEL org.opencontainers.image.source="https://github.com/UniPro-tech/UniBot"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="UniBot"
LABEL org.opencontainers.image.description="A multifunctional Discord bot for community management, entertainment, and productivity."
LABEL org.opencontainers.image.vendor="All-Japan Digital Creative Club UniProject <info@uniproject.jp>"

WORKDIR /app

RUN apt-get update && apt-get install -y \
    curl git file unzip \
    python3 make g++ ffmpeg \
    && rm -rf /var/lib/apt/lists/*

# bun だけ入れる（Nodeは20のまま）
RUN curl -fsSL https://bun.sh/install | bash
ENV PATH="/root/.bun/bin:$PATH"

COPY package*.json ./
COPY bun.lockb ./

RUN bun install --trust-all

COPY . .

CMD ["sh", "-c", "bunx prisma db push && bun run src/index.ts"]


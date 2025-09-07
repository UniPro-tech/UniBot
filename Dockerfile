FROM oven/bun:slim
WORKDIR /app

LABEL org.opencontainers.image.authors="Yuito Akatsuki <yuito@yuito-it.jp>"
LABEL org.opencontainers.image.url="https://github.com/UniPro-tech/UniBot"
LABEL org.opencontainers.image.source="https://github.com/UniPro-tech/UniBot"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="UniBot"
LABEL org.opencontainers.image.description="A multifunctional Discord bot for community management, entertainment, and productivity."
LABEL org.opencontainers.image.vendor="All-Japan Digital Creative Club UniProject <info@uniproject.jp>"

RUN apt-get update && apt-get install -y \
    curl \
    git \
    file \
    unzip

RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | tee /usr/share/keyrings/github-archive-keyring.gpg >/dev/null && \
    echo "deb [signed-by=/usr/share/keyrings/github-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list && \
    apt-get update && \
    apt-get install -y gh

RUN apt-get update && apt-get install -y python3 make g++ ffmpeg && rm -rf /var/lib/apt/lists/*

COPY package*.json ./
COPY bun.lock ./
RUN bun install --trust-all

COPY . .

CMD ["sh", "-c", "bunx prisma db push && bun run src/index.ts"]

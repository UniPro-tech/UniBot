FROM oven/bun:slim
WORKDIR /app

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
RUN bun install

COPY . .

CMD ["sh", "-c", "bunx prisma db push && bun run src/index.ts"]

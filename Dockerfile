FROM oven/bun:slim
WORKDIR /app

COPY package*.json ./
RUN bun install

COPY . .

CMD [ "bun", "run", "src/index.ts" ]

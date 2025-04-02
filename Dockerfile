FROM node:slim
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

CMD [ "bun", "run", "src/index.ts" ]

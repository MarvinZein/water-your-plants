FROM node:21-alpine

WORKDIR /app

COPY . .

RUN npm ci

RUN npm run build

RUN rm -rf src/ static/ docker-compose.yml

USER node:node

CMD ["node","build/index.js"]
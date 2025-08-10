FROM node:21-alpine3.20
LABEL Maintainer="stepan-vinokurov"


WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build

# Устанавливаем serve для обслуживания статических файлов
RUN npm install -g serve

# Запускаем сервер
CMD ["serve", "-s", "dist", "-l", "8080"]

EXPOSE 8080

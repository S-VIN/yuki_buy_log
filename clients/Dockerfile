FROM node:21-alpine3.20
LABEL Maintainer="stepan-vinokurov"


WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .
RUN npm run build

# Install serve to serve static files
RUN npm install -g serve

# Start the server
CMD ["serve", "-s", "dist", "-l", "8080"]

EXPOSE 8080

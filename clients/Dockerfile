FROM node:21-alpine3.20 AS builder
LABEL Maintainer="stepan-vinokurov"

WORKDIR /app

# Copy package files
COPY package.json ./

# Install dependencies with production optimizations
RUN npm install --production=false --no-audit

# Copy source code
COPY . .

# Build argument for API URL (can be overridden at build time)
ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}

# Build the application with production optimizations
RUN npm run build

# Production stage
FROM node:21-alpine3.20
WORKDIR /app

# Install serve for serving static files
RUN npm install -g serve@14.2.3

# Copy built files from builder
COPY --from=builder /app/dist ./dist

EXPOSE 8080

# Start the server
CMD ["serve", "-s", "dist", "-l", "8080"]

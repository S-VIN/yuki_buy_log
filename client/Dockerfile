FROM node:latest AS builder
LABEL Maintainer="stepan-vinokurov"

WORKDIR /app

# Copy package files
COPY package.json package-lock.json ./

# Install dependencies with production optimizations
RUN npm ci --no-audit

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

EXPOSE 8000

# Start the server - listen on 0.0.0.0 to accept external connections
CMD ["serve", "-s", "dist", "-l", "tcp://0.0.0.0:8000"]
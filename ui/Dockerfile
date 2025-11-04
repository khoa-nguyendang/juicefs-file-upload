# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Set build-time environment variables (these are needed during build)
ARG NEXT_PUBLIC_API_URL=http://go-app:8080
ARG NEXT_PUBLIC_API_TIMEOUT=30000
ARG NEXT_PUBLIC_MINIO_URL=http://minio:9000
ARG NEXT_PUBLIC_MINIO_CONSOLE_URL=http://localhost:9001
ARG NEXT_PUBLIC_APP_NAME="JuiceFS File Browser"
ARG NEXT_PUBLIC_MAX_UPLOAD_SIZE=104857600
ARG NEXT_PUBLIC_ENABLE_DELETE=true
ARG NEXT_PUBLIC_ENABLE_UPLOAD=true
ARG NEXT_PUBLIC_DEFAULT_VIEW=tree
ARG NEXT_PUBLIC_SHOW_HIDDEN_FILES=false
ARG NEXT_PUBLIC_FILE_PREVIEW_ENABLED=true

# Build the application
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Copy built application from builder
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

# Install production dependencies only
# RUN npm ci --only=production

# Environment variables (can be overridden at runtime)
ENV NEXT_PUBLIC_API_URL=http://go-app:8080
ENV NEXT_PUBLIC_API_TIMEOUT=30000
ENV NEXT_PUBLIC_MINIO_URL=http://minio:9000
ENV NEXT_PUBLIC_MINIO_CONSOLE_URL=http://localhost:9001
ENV NEXT_PUBLIC_APP_NAME="JuiceFS File Browser"
ENV NEXT_PUBLIC_MAX_UPLOAD_SIZE=104857600
ENV NEXT_PUBLIC_ENABLE_DELETE=true
ENV NEXT_PUBLIC_ENABLE_UPLOAD=true
ENV NEXT_PUBLIC_DEFAULT_VIEW=tree
ENV NEXT_PUBLIC_SHOW_HIDDEN_FILES=false
ENV NEXT_PUBLIC_FILE_PREVIEW_ENABLED=true

# Expose port
EXPOSE 3000

# Start the application
CMD ["node", "server.js"]
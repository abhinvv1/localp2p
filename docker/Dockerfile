FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /app

# Copy Go modules
COPY core/go.mod core/go.sum ./
RUN go mod download

# Copy source code
COPY core/ ./

# Build the application
RUN go build -o localp2p

FROM node:18-alpine

# Install system dependencies
RUN apk add --no-cache curl

# Set working directory
WORKDIR /app

# Copy built Go binary
COPY --from=builder /app/localp2p ./core/localp2p

# Copy Node.js CLI
COPY cli/ ./cli/

# Install CLI dependencies
WORKDIR /app/cli
RUN npm install

# Make CLI executable
RUN chmod +x bin/localp2p

# Expose ports
EXPOSE 8080 9090

# Set working directory back to app root
WORKDIR /app

# Default command
CMD ["./cli/bin/localp2p", "start"]
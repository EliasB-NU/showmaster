# Build the Go application
FROM golang:1.24.1-alpine AS builder-go

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY src/ ./src/

RUN go build ./src/main.go

# Build the Node.js application
FROM node:22.13.0-alpine AS builder-node

WORKDIR /app

COPY web/ ./

RUN npm install

RUN npm run build

# Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

COPY --from=builder-go /app/main .

COPY templates/ ./templates/

COPY --from=builder-node /app/dist ./web/dist

EXPOSE 3001 3002

CMD ["./main", "prod"]
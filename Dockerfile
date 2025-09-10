# Stage 1: Build the Go applications
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the API server
RUN CGO_ENABLED=0 go build -o /sistem-manajemen-armada-api ./api

# Build the MQTT publisher
RUN CGO_ENABLED=0 go build -o /sistem-manajemen-armada-publisher ./mqtt-publisher

# Build the RabbitMQ worker
RUN CGO_ENABLED=0 go build -o /sistem-manajemen-armada-worker ./rabbitmq-worker

# Stage 2: Create a minimal image for the backend service
FROM alpine:latest AS backend

WORKDIR /app

COPY --from=builder /sistem-manajemen-armada-api .

EXPOSE 9999

CMD ["./sistem-manajemen-armada-api"]

# Stage 3: Create a minimal image for the publisher service
FROM alpine:latest AS publisher

WORKDIR /app

COPY --from=builder /sistem-manajemen-armada-publisher .

CMD ["./sistem-manajemen-armada-publisher"]

# Stage 4: Create a minimal image for the worker service
FROM alpine:latest AS worker

WORKDIR /app

COPY --from=builder /sistem-manajemen-armada-worker .

CMD ["./sistem-manajemen-armada-worker"]
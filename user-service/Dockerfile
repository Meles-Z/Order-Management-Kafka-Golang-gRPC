#syntax=docker/dockerfile:1
FROM golang:1.24-alpine

# Install system dependencies and librdkafka
RUN apk --no-cache update && \
apk --no-cache add git gcc libc-dev

# Kafka Go client is based on the C library librdkafka
ENV CGO_ENABLED=1

WORKDIR /app

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# Copy Go module files first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Now copy the entire project
COPY . .

# Expose the port your app will run on
EXPOSE 8080

# Run Air (hot reload)
CMD ["air", "-c", ".air.toml"]

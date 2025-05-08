
#syntax=docker/dockerfile:1
FROM golang:1.24-alpine

# Install system dependencies
RUN apk add --no-cache git curl

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

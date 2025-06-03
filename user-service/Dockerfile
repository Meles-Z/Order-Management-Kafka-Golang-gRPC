FROM golang:1.24-alpine

WORKDIR /app

# Install Air (for hot reload)
RUN go install github.com/air-verse/air@latest

# Copy Go mod files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# Start Air (will watch for changes)
CMD ["air", "-c", ".air.toml"]

FROM golang:1.24-alpine

# Install system dependencies
RUN apk add --no-cache git curl build-base librdkafka-dev

WORKDIR /app

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

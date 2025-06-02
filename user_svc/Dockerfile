FROM golang:1.24

WORKDIR /app

# Install build tools and librdkafka
RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    gcc \
    make \
    git \
    && rm -rf /var/lib/apt/lists/*

# Install Air
RUN go install github.com/air-verse/air@latest

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy code
COPY . .

# Start app with Air
CMD ["air", "-c", ".air.toml"]
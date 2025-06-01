FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct
RUN echo "precedence ::ffff:0:0/96  100" >> /etc/gai.conf

RUN go mod download

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]

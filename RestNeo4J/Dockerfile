FROM golang:alpine AS builder
WORKDIR /app
COPY . .

EXPOSE 89
ENTRYPOINT ["go", "run", "main.go"]
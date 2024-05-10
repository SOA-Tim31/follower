FROM golang:alpine AS builder
WORKDIR /app
COPY . .

EXPOSE 8083
ENTRYPOINT ["go", "run", "main.go"]




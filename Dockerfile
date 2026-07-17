# Stage 1: build
FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scheduler .

# Stage 2: runtime
FROM ubuntu:24.04

WORKDIR /app

COPY --from=builder /build/scheduler .
COPY --from=builder /build/web       ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=

EXPOSE 7540

CMD ["./scheduler"]

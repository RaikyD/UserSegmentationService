# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/UserSegmentationService
RUN go build -o /usr/bin/app

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /usr/bin/app /usr/bin/app
COPY migrations /app/migrations
ENV HTTP_PORT=8080
EXPOSE 8080
CMD ["/usr/bin/app"]
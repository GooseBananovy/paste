FROM golang:1.26.1-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o paste ./cmd/server

FROM alpine:latest
COPY --from=builder /app/paste /paste
CMD ["/paste"]

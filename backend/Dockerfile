FROM golang:alpine AS builder

WORKDIR /app
ENV GOTOOLCHAIN=auto
RUN apk add --no-cache git ca-certificates


COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server


FROM alpine:latest

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && mkdir -p /app/data

COPY --from=builder /app/server .
COPY --from=builder /app/.env ./.env

EXPOSE 8080

CMD ["./server"]

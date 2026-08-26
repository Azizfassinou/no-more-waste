# Multi-stage Dockerfile pour le backend Go (No More Waste)

# Étape 1 : Build du binaire Go
FROM golang:alpine AS builder

WORKDIR /app
ENV GOTOOLCHAIN=auto

# Installation des dépendances système
RUN apk add --no-cache git ca-certificates

# Cache des dépendances Go
COPY go.mod go.sum ./
RUN go mod download

# Copie du code source
COPY . .

# Compilation des binaires statiques ultralégers
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o seed seeder.go

# Étape 2 : Image d'exécution minimale
FROM alpine:latest

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && mkdir -p /app/data

# Copie des binaires et de l'environnement
COPY --from=builder /app/server .
COPY --from=builder /app/seed .
COPY --from=builder /app/.env ./.env

EXPOSE 8080

CMD ["./server"]

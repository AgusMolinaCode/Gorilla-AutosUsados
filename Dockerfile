# Build stage
FROM golang:1.22-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache gcc musl-dev

WORKDIR /build

# Copia go.mod y go.sum
COPY go.mod go.sum ./

# Descarga las dependencias
RUN go mod download

# Copia el código fuente
COPY . .

# Compila la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copia el ejecutable y el archivo .env
COPY --from=builder /build/app .
COPY .env .

# Expone el puerto
EXPOSE 8080

# Ejecuta la aplicación
CMD ["./app"]

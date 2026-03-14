# ============================================================
# Stage 1: Builder
# ============================================================
FROM golang:1.23-alpine AS builder

# Instalar dependencias del sistema necesarias para CGO (sqlite y otros drivers)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copiar los archivos de dependencias primero para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar el binario del servidor
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server/main.go

# ============================================================
# Stage 2: Runner (imagen final mínima)
# ============================================================
FROM alpine:3.19

# Instalar certificados TLS y librerías C necesarias en tiempo de ejecución
RUN apk add --no-cache ca-certificates tzdata libc6-compat

WORKDIR /app

# Copiar el binario compilado desde el stage builder
COPY --from=builder /app/server .

# Crear directorio para uploads (archivos subidos por usuarios)
RUN mkdir -p /app/uploads

# Exponer el puerto de la aplicación
EXPOSE 8080

# Ejecutar el servidor
CMD ["./server"]

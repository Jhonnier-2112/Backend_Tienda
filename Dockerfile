# ============================================================
# Stage 1: Builder
# ============================================================
FROM golang:1.25-alpine AS builder

# Dependencias necesarias para CGO
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copiar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código
COPY . .

# Compilar binario
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server ./cmd/server/main.go


# ============================================================
# Stage 2: Runner
# ============================================================
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata libc6-compat

WORKDIR /app

# Copiar binario
COPY --from=builder /app/server .

# Crear carpeta uploads
RUN mkdir -p /app/uploads

# Railway usa puerto dinámico
ENV PORT=10000

EXPOSE 10000

CMD ["./server"]
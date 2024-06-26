# Etapa de construcción
FROM golang:1.22-alpine AS builder

# Instalar Git para que 'go get' funcione
RUN apk update && apk add --no-cache git

WORKDIR /go/src/app

COPY . .

# Descargar las dependencias y construir el binario
RUN go mod tidy
RUN go build -o /go/bin/app

# Etapa final
FROM alpine:latest

# Instalar certificados CA
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /go/bin/app /app

ENTRYPOINT ["/app/app"]

EXPOSE 3000

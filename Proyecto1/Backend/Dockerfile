FROM golang:alpine AS builder

WORKDIR /go/src/app

# Copiar el contenido de Backend
COPY . ./

RUN apk add --no-cache git
RUN go get -d -v
RUN go build -o /go/bin/app
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get go.mongodb.org/mongo-driver/mongo/options
RUN go get github.com/gofiber/fiber/v2
RUN go get github.com/gofiber/fiber/v2/middleware/cors

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash sysstat procps

WORKDIR /app
COPY --from=builder /go/bin/app /app
COPY .env /app/.env

ENTRYPOINT ["/app/app"]

EXPOSE 8080

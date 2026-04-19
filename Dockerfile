FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o orders-app ./cmd/server

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates curl

COPY --from=builder /app/orders-app /app/orders-app
COPY --from=builder /app/web /app/web

EXPOSE 8080

CMD ["/app/orders-app"]

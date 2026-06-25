FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/stops-api ./cmd/api
RUN go build -o /out/stops-migrate ./cmd/migrate

FROM alpine:3.21

WORKDIR /app

RUN adduser -D -g "" appuser
COPY --from=builder /out/stops-api /usr/local/bin/stops-api
COPY --from=builder /out/stops-migrate /usr/local/bin/stops-migrate

USER appuser

EXPOSE 8080

CMD ["stops-api"]

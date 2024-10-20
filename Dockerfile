FROM golang:1.23-alpine AS builder

COPY . /eventflow/source/
WORKDIR /eventflow/source/

RUN go build -o ./bin/eventflow cmd/api/main.go
RUN go build -o ./bin/migrator cmd/migrator/main.go

FROM alpine:3.13

WORKDIR /root/
COPY --from=builder /eventflow/source/bin/ .
COPY --from=builder /eventflow/source/migrations /migrations/
COPY --from=builder /eventflow/source/prod.env .

CMD ["sh", "-c", "./migrator --env-path=prod.env --migrations-path=/migrations/ && ./eventflow --env-path=prod.env --env-level=prod --logs-path=/eventflow/logs"]
FROM golang:1.22.5-alpine AS builder

COPY . /eventflow/source/
WORKDIR /eventflow/source/

RUN go mod download
RUN go build -o ./bin/migrator cmd/migrator/main.go

FROM alpine:3.13

WORKDIR /root/
COPY --from=builder /eventflow/source/bin/ .
COPY --from=builder /eventflow/source/migrations /migrations/
COPY --from=builder /eventflow/source/prod.env .

CMD ["./migrator", "--env-path=prod.env", "--migrations-path=/migrations/"]
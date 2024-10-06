FROM golang:1.22.5-alpine AS builder

COPY . /eventflow/source/
WORKDIR /eventflow/source/

RUN go build -o ./bin/eventflow cmd/api/main.go

FROM alpine:3.13

WORKDIR /root/
COPY --from=builder /eventflow/source/bin/eventflow .
COPY --from=builder /eventflow/source/prod.env .

CMD ["./eventflow", "--env-path=prod.env", "--env-level=prod", "--logs-path=/eventflow/logs"]
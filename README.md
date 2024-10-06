# EventFlow backend

EventFlow is a self-hosted platform for managing and selling tickets to events. 
This backend service is built with Go and provides a REST API 
for event management, ticket sales, and payment processing. 
The platform integrates with Yookassa for handling payments 
and supports sending emails for event notifications. 
For observability, EventFlow includes monitoring and logging 
through Grafana, Prometheus, and Grafana Loki. 
The repository also includes tests for api and service layer.

**If you want to test payments, use test card from [here](https://yookassa.ru/developers/payment-acceptance/testing-and-going-live/testing?lang=en)**

## Features
* Events management;
* Yookassa payments integration;
* Email notifications;
* Observability;
* Tests.

## How to run
1. Configure .env according to .env.example
2. Clone repo
   ```shell
   git clone https://github.com/wDRxxx/eventflow-backend.git
   ```
3. Run database, mailhog and apply migrations
   ```shell
   docker-compose up -d
   make migrations-up
   ```

4. Run observability services (optionally)
   ```shell
   docker-compose -f=./observability/docker-compose.yaml up -d
   ```

5. Run backend
   ```shell
   go run ./cmd/api/main.go
   ```

6. Run [frontend](https://github.com/wDRxxx/eventflow-frontend) (if you want)

### Tests
To run tests use:
```shell
make test
```

If you want to see coverage, run:
```shell
make test-coverage
```
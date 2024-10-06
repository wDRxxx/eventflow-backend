migrations-up:
	go run ./cmd/migrator/main.go \
	--env-path=.env \
	--migrations-path=migrations

migrations-down:
	go run ./cmd/migrator/main.go \
	--action=down \
	--env-path=./.env \
	--migrations-path=./migrations

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/wDRxxx/eventflow-backend/internal/api/...,github.com/wDRxxx/eventflow-backend/internal/service/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/wDRxxx/eventflow-backend/internal/api/...,github.com/wDRxxx/eventflow-backend/internal/service/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=./coverage.out | grep "total";
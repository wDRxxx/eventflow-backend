migrations-up:
	go run ./cmd/migrator/main.go \
	--env-path=./.env \
	--migrations-path=./migrations

migrations-down:
	go run ./cmd/migrator/main.go \
	--action=down \
	--env-path=./.env \
	--migrations-path=./migrations
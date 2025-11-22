BUF_VERSION = v1.45.0

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test.cleancache
test.cleancache:
	go clean -testcache

.PHONY: test.unit
test.unit: test.cleancache
	go test -race ./...

.PHONY: test.cover
test.cover: test.cleancache
	go test -v -race ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out

.PHONY: migration
migration:
	migrate create -ext sql -dir migration -seq ${name}

.PHONY: migrate
migrate:
	go run cmd/migrate/main.go 

.PHONY: rollback
rollback:
	go run cmd/migrate/main.go -down

.PHONY: seed
seed:
	go run cmd/seed/main.go
	
.PHONY: run
run:
	go run cmd/rest_api/main.go

.PHONY: run.dev
run.dev:
	air
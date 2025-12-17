tidy ::
	@go mod tidy

seed ::
	@go run cmd/seed/main.go

run ::
	@sed -i 's/^POSTGRES_HOST=challenge-database$$/POSTGRES_HOST=localhost/' .env
	@go run cmd/server/main.go

test ::
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

coverage: test
	@go tool cover -html=coverage.out -o coverage.html
	@firefox coverage.html

docker-up ::
	docker compose up -d postgres

docker-down ::
	docker compose down

server ::
	@sed -i 's/^POSTGRES_HOST=localhost$$/POSTGRES_HOST=challenge-database/' .env
	docker compose up

tidy ::
	@go mod tidy

seed ::
	@sed -i 's/^POSTGRES_HOST=challenge-database$$/POSTGRES_HOST=localhost/' .env
	@go run cmd/seed/main.go

run ::
	@sed -i 's/^POSTGRES_HOST=challenge-database$$/POSTGRES_HOST=localhost/' .env
	@go run cmd/server/main.go

test ::
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

coverage: test
	@grep -vFf .covignore coverage.out > coverage.filtered.out
	@go tool cover -html=coverage.filtered.out -o coverage.html
	@rm -f coverage.filtered.out
	@command -v xdg-open >/dev/null && xdg-open coverage.html || open coverage.html

docker-up ::
	@sed -i 's/^POSTGRES_HOST=localhost$$/POSTGRES_HOST=challenge-database/' .env
	@docker compose up

docker-up-db ::
	@docker compose up -d postgres

docker-down ::
	docker compose down

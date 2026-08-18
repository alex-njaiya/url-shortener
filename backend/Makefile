run:
	@go run cmd/api/main.go

migrate-up:
	@migrate -path internal/database/migrations -database "postgres://alex:41484877@localhost:5432/url_shortener?sslmode=disable" up

migrate-down:
	@migrate -path internal/database/migrations -database "postgres://alex:41484877@localhost:5432/url_shortener?sslmode=disable" down
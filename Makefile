clean-up:
	@docker compose down -v
	@docker compose up --build

run-backend:
	@docker compose up -d backend

migrate-up:
	@migrate -path backend/internal/database/migrations -database "postgres://alex:41484877@localhost:5432/url_shortener?sslmode=disable" up

migrate-down:
	@migrate -path backend/internal/database/migrations -database "postgres://alex:41484877@localhost:5432/url_shortener?sslmode=disable" down

compose-up:
	@docker compose up --build 
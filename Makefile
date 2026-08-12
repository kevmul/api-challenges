include .env
export

migrate-up:
	migrate -path ./migrations -database $(DATABASE_URL) up

migrate-down:
	migrate -path ./migrations -database $(DATABASE_URL) down

migrate-drop:
	migrate -path ./migrations -database $(DATABASE_URL) drop -f

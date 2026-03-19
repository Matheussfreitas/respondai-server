include .env
export

migrate-up:
	migrate -database "$(DATABASE_URL)" -path migrations up

migrate-down:
	migrate -database "$(DATABASE_URL)" -path migrations down

.PHONY: swagger
swagger:
	GOCACHE=$${GOCACHE:-/tmp/go-build-cache} go generate ./cmd/api

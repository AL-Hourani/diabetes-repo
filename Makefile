build:
	@go build -o bin/care-center main.go

test:
	@go test -v ./...


run: build 
	@./bin/care-center


migration:
	@migrate create -ext sql -dir migrate/migrations $(filter-out $@, $(MAKECMDGOALS))

migrate-up:
	@go run migrate/main.go up

migrate-down:
	@go run migrate/main.go down
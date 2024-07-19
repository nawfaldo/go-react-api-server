build:
	@go build -o bin/test cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/test


migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -path cmd/migrate/migrations -database "mysql://root@tcp(localhost:3306)/go_test" up

migrate-down:
	@migrate -path cmd/migrate/migrations -database "mysql://root@tcp(localhost:3306)/go_test" down

migrate-force:
	@migrate -path cmd/migrate/migrations -database "mysql://root@tcp(localhost:3306)/go_test" force $(filter-out $@,$(MAKECMDGOALS))

%:
	@:

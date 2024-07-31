DBUSER=root
DBPASSWORD=password
DBPORT=5432
DBNAME=fraud_detection
DSN=postgresql://$(DBUSER):$(DBPASSWORD)@localhost:$(DBPORT)/$(DBNAME)?sslmode=disable

db-up:
	docker run --name postgres -e POSTGRES_USER=$(DBUSER) -e POSTGRES_PASSWORD=$(DBPASSWORD) -p $(DBPORT):$(DBPORT) -d postgres

db-down:
	@if docker ps -a --format '{{.Names}}' | grep -q '^postgres$$'; then \
        docker stop postgres; \
        docker rm postgres; \
        echo "Container stopped."; \
    else \
        echo "Container not found, skipping stop command."; \
    fi

db-connect:
	docker exec -it postgres psql -U $(DBUSER) -p $(DBPORT) -d postgres

db-migrate:
	docker exec -it postgres psql -U $(DBUSER) -p $(DBPORT) -d postgres -exec "CREATE DATABASE $(DBNAME);"
	goose -dir=./migrations postgres "$(DSN)" up

db-reset:
	docker exec -it postgres psql -U $(DBUSER) -p $(DBPORT) -d postgres -exec "DROP DATABASE $(DBNAME);"

db-drop:
	docker exec -it postgres psql -U $(DBUSER) -p $(DBPORT) -d postgres -exec "DROP DATABASE $(DBNAME);"

run:
	go run consumers/cmd/main.go

build:
	CGO_ENABLED=0 go build -o ./service-catalog-api app/cmd/main.go

test:
	go test -v ./...

format:
	go fmt ./...

spec-up:
	docker run -d --rm --name swagger -p 8080:8080 -v $$PWD/spec:/spec -e SWAGGER_JSON=/spec/openapi.yaml swaggerapi/swagger-ui

spec-down:
	docker stop swagger
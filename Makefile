all: build

build:
	@echo "Building..."
	@make build-server
	@make build-client

build-server:
	@echo "Building server..."
	@go build -o bin/aggregation-server aggregation-server/cmd/main.go

build-client:
	@echo "Building client..."
	@go build -o bin/he-client he-client/cmd/main.go

run-server:
	@echo "Running server..."
	@go run aggregation-server/cmd/main.go

run-client:
	@echo "Running client..."
	@go run he-client/cmd/main.go

docker-run:
	@echo "Running docker..."
	@if docker compose up -f deployment/docker-compose.yml --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build -f deployment/docker-compose.yml; \
	fi

docker-down:
	@if docker compose down -f deployment/docker-compose.yml 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down -f deployment/docker-compose.yml; \
	fi

clean:
	@echo "Cleaning..."
	@rm -rf bin

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi
	
.PHONY: all build run-server run-client clean watch docker-run docker-down

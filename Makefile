define setup_env
	$(eval ENV_FILE := env/$(1).env)
	@echo " - setup env $(ENV_FILE)"
	$(eval include env/$(1).env)
	$(eval export sed 's/=.*//' env/$(1).env)
endef




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

server-run:
	@echo "Running server..."
	@$(call setup_env,server)
	@go run aggregation-server/cmd/main.go

client-run:
	@echo "Running client..."
	@$(call setup_env,client)
	@go run he-client/cmd/main.go

python-client-run:
	@echo "Running python client..."
	@$(call setup_env,client)
	@python3 python-client/src/main.py

python-generate-image-data:
	@python3 -c "import numpy as np; import random; print([[random.random() for x in range(28)] for y in range(28)])"

docker-run:
	@echo "Running docker..."
	@if docker compose -f deployment/docker-compose.yml up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f deployment/docker-compose.yml up --build; \
	fi

docker-down:
	@if docker compose -f deployment/docker-compose.yml down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f deployment/docker-compose.yml down; \
	fi

podman-build:
	podman build -t he-aggregation-server:latest ./aggregation-server/
	podman build -t he-client:latest ./he-client/
	podman build -t he-python-client:latest ./python-client/

podman-run:
	podman network create --ignore dshe
	podman play kube --network dshe --replace ./deployment/kube/aggregation-server.yaml
	podman play kube --network dshe --replace ./deployment/kube/he-client-1.yaml

podman-down:
	podman play kube --network dshe --down ./deployment/kube/aggregation-server.yaml
	podman play kube --network dshe --down ./deployment/kube/he-client-1.yaml
	podman network rm dshe

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

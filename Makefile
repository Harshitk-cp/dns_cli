# Makefile

CLI_BINARY_NAME=dnscli

all: build

build:
	@echo "Building the CLI..."
	go build -o $(CLI_BINARY_NAME) ./cmd/query-tool

alias:
	@echo "Creating an alias..."
	sudo ln -sf $(PWD)/$(CLI_BINARY_NAME) /usr/local/bin/$(CLI_BINARY_NAME)

docker-build:
	@echo "Building the DNS server Docker image..."
	docker-compose build

run-server:
	@echo "Running the DNS and DoH servers..."
	docker-compose up

stop-server:
	@echo "Stopping the DNS and DoH servers..."
	docker-compose down

clean:
	@echo "Cleaning up..."
	rm -f $(CLI_BINARY_NAME)
	sudo rm -f /usr/local/bin/$(CLI_BINARY_NAME)

.PHONY: all build alias docker-build run-server stop-server clean

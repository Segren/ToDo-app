include .env

# ==================================================================================== # 
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]


COLOR_GREEN = "\e[1;32m%s\e[0m\n"
SWAGGER_VERSION		= v1.16.2

# ==================================================================================== # 
# DEVELOPMENT
# ==================================================================================== #

## cmd/toDo-app: run the cmd/toDo-app application
.PHONY: cmd/toDo-app
cmd/toDo-app: 
	go run ./cmd/toDo-app/

## docker/up: launch docker
.PHONY: docker/up
docker/up: confirm
	docker-compose up 

## docker/down: terminate docker
.PHONY: docker/down
docker/down: confirm
	docker-compose down

## swagger/install: install the swag cli
.PHONY: swagger/install
swagger/install:
	@printf $(COLOR_GREEN) "Install/update swag cli"
	go install github.com/swaggo/swag/cmd/swag@$(SWAGGER_VERSION)
	@printf $(COLOR_GREEN) "Success [install swaggo]"

## swagger/init: Generates the godoc !(Requires swagger installed)
.PHONY: swagger/init
swagger/init:
	@printf $(COLOR_GREEN) "Generating swagger documentation\n"
	@cd $(shell git rev-parse --show-toplevel) && swag init --parseDependency --parseInternal --quiet --generatedTime -g cmd/toDo-app/main.go -d .
	@printf $(COLOR_GREEN) "Success [gen swagger]\n"

# ==================================================================================== # 
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit 
audit: vendor
	@echo 'Tidying and verifying module dependencies...' go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== # 
# BUILD
# ==================================================================================== #

current_time = $(shell date "+%Y-%m-%dT%H:%M:%S%z")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build/toDo-app
build/toDo-app: 
	@echo 'Building cmd/toDo-app...'
	go build -ldflags=${linker_flags} -o=./bin/toDo-app ./cmd/toDo-app
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/toDo-app ./cmd/toDo-app
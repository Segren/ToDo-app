# Colors
COLOR_GREEN = "\e[1;32m%s\e[0m\n"

# Utils versions variables
SWAGGER_VERSION		= v1.16.2

all: 
	go run cmd/toDo-app/main.go -app-config=.env

up: 
	docker-compose up 

down: 
	docker-compose down

install-swagger: ## Install the swag cli
	@printf $(COLOR_GREEN) "Install/update swag cli"
	go install github.com/swaggo/swag/cmd/swag@$(SWAGGER_VERSION)
	@printf $(COLOR_GREEN) "Success [install swaggo]"

swagger-init: ## Generates the godoc !(Requires swagger installed)
	@printf $(COLOR_GREEN) "Generating swagger documentation\n"
	@cd $(shell git rev-parse --show-toplevel) && swag init --parseDependency --parseInternal --quiet --generatedTime -g cmd/toDo-app/main.go -d .
	@printf $(COLOR_GREEN) "Success [gen swagger]\n"

si: swagger-init
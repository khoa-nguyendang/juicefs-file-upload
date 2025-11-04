.PHONY: build up down build-and-push

# Variables
DOCKER_COMPOSE = docker-compose
GO = go
CURL = curl

# Colors for output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[1;33m
NC = \033[0m # No Color



build: ## Build all Docker images
	@echo "${YELLOW}Building Docker images...${NC}"
	@$(DOCKER_COMPOSE) build
	@echo "${GREEN}Build complete${NC}"

up: ## Start all services
	@echo "${YELLOW}Starting all services...${NC}"
	@$(DOCKER_COMPOSE) up -d
	@echo "${GREEN}Services started${NC}"

down: ## Stop all services
	@echo "${YELLOW}Stopping all services...${NC}"
	@$(DOCKER_COMPOSE) down
	@echo "${GREEN}Services stopped${NC}"

build-and-push: ## Build all and push images
	@docker build -t naturemyloves/file-browser-server:latest ./server
	@docker build -t naturemyloves/file-browser-ui:latest ./ui
	@docker build -t naturemyloves/displaying-sharing:latest ./display-sharing
	@docker push naturemyloves/file-browser-server:latest
	@docker push naturemyloves/file-browser-ui:latest
	@docker push naturemyloves/displaying-sharing:latest
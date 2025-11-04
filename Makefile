.PHONY: build up down

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

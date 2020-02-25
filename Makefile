.PHONY: start
start:
	docker-compose -f deployments/docker-compose.yaml up --build todo-backend

.PHONY: specs
specs:
	docker-compose -f deployments/docker-compose.yaml up todo-backend-specs

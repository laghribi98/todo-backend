.PHONY: start
start:
	docker-compose -f deployments/docker-compose.yaml up --build todo-backend

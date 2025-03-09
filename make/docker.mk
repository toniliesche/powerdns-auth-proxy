include make/docker_compose.mk
include make/docker_database.mk

create-network:
	@if ! docker network ls --format '{{.Name}}' | grep -q "^pdns-network$$"; then \
		docker network create pdns-network; \
		echo "Network pdns-network created."; \
	else \
		echo "Network pdns-network already exists."; \
	fi

AUTH_PROXY_CMD=docker compose -f docker/docker-compose.yml -p pdns-auth-proxy
MARIADB_CMD=docker compose -f docker/docker-compose.mariadb.yml -p pdns-database

up: up-mariadb setup-mariadb-powerdns setup-mariadb-gateway up-gateway-mariadb;

up-gateway: create-network
	@docker compose -p powerdns-auth-proxy up -d

up-gateway-jwt: create-network
	@AUTH_TYPE=jwt $(AUTH_PROXY_CMD) up -d --remove-orphans

up-gateway-mariadb: create-network up-mariadb check-db
	@AUTH_TYPE=jwt DB_TYPE=mysql $(AUTH_PROXY_CMD) up -d --remove-orphans

up-mariadb: create-network
	@$(MARIADB_CMD) up --pull always -d --remove-orphans

start: start-gateway start-mariadb;

start-gateway:
	@$(AUTH_PROXY_CMD) start

start-mariadb:
	@$(MARIADB_CMD) start

down: down-gateway down-mariadb;

down-gateway:
	@$(AUTH_PROXY_CMD) down --volumes

down-mariadb:
	@$(MARIADB_CMD) down --volumes

stop: stop-gateway stop-mariadb;

stop-gateway:
	@$(AUTH_PROXY_CMD) stop

stop-mariadb:
	@$(MARIADB_CMD) stop

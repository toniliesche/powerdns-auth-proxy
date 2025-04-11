# MIT License
# Copyright (c) 2025 Toni Liesche
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

AUTH_PROXY_CMD=docker compose -f docker/docker-compose.yml -p powerdns-auth-proxy
MARIADB_CMD=docker compose -f docker/docker-compose.mariadb.yml -p powerdns-mariadb
POSTGRES_CMD=docker compose -f docker/docker-compose.postgres.yml -p powerdns-postgres

up: up-mariadb setup-mariadb-powerdns setup-mariadb-gateway up-gateway-mariadb;

up-pg: up-mariadb up-postgres setup-mariadb-powerdns setup-postgres-gateway up-gateway-postgres;

up-gateway: create-network
	@docker compose -p powerdns-auth-proxy up -d

up-gateway-jwt: create-network
	@AUTH_TYPE=jwt $(AUTH_PROXY_CMD) up -d --remove-orphans

up-gateway-mariadb: create-network up-mariadb check-db-mariadb
	@AUTH_TYPE=jwt DB_TYPE=mysql $(AUTH_PROXY_CMD) up -d --remove-orphans

up-gateway-postgres: create-network up-postgres check-db-postgres
	@AUTH_TYPE=jwt DB_TYPE=postgres $(AUTH_PROXY_CMD) up -d --remove-orphans

up-mariadb: create-network
	@$(MARIADB_CMD) up --pull always -d --remove-orphans

up-postgres: create-network
	@$(POSTGRES_CMD) up --pull always -d --remove-orphans

start: start-gateway start-mariadb;

start-gateway:
	@$(AUTH_PROXY_CMD) start

start-mariadb:
	@$(MARIADB_CMD) start

start-postgres:
	@$(POSTGRES_CMD) start

down: down-gateway down-mariadb down-postgres;

down-gateway:
	@$(AUTH_PROXY_CMD) down --volumes

down-mariadb:
	@$(MARIADB_CMD) down --volumes

down-postgres:
	@$(POSTGRES_CMD) down --volumes

stop: stop-gateway stop-mariadb stop-postgres;

stop-gateway:
	@$(AUTH_PROXY_CMD) stop

stop-mariadb:
	@$(MARIADB_CMD) stop

stop-postgres:
	@$(POSTGRES_CMD) stop

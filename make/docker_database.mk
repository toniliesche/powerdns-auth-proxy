MARIADB_CMD=docker compose -f docker/docker-compose.mariadb.yml -p pdns-mariadb
POSTGRES_CMD=docker compose -f docker/docker-compose.postgres.yml -p pdns-postgres

setup-mariadb-gateway: up-mariadb check-db-mariadb create-gateway-db-mariadb create-gateway-db-mariadb-user

setup-mariadb-powerdns: up-mariadb check-db-mariadb create-powerdns-db-mariadb create-powerdns-db-mariadb-user

setup-mariadb: setup-mariadb-gateway setup-mariadb-powerdns

setup-postgres-gateway: up-postgres check-db-postgres create-gateway-db-postgres create-gateway-db-postgres-user

setup-postgres-powerdns: up-postgres check-db-postgres create-powerdns-db-postgres create-powerdns-db-postgres-user

setup-postgres: setup-postgres-gateway setup-postgres-powerdns

check-db-mariadb: up-mariadb
	@echo "Checking if MariaDB is available..."
	@until $(MARIADB_CMD) exec mariadb mariadb-admin -u $(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) ping > /dev/null 2>&1; do \
		echo "Waiting for MariaDB to accept connections..."; \
		sleep 1; \
	done
	@echo "MariaDB is up and accepting connections."
	@sleep 5;

create-gateway-db-mariadb: check-db-mariadb
	@if ! $(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "SHOW DATABASES LIKE '$(DB_MARIADB_GATEWAY_DATABASE)';" | grep -q "$(DB_MARIADB_GATEWAY_DATABASE)"; then \
		echo "Creating database $(DB_MARIADB_GATEWAY_DATABASE)..."; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "CREATE DATABASE $(DB_MARIADB_GATEWAY_DATABASE);"; \
		echo "Database $(DB_MARIADB_GATEWAY_DATABASE) created."; \
	else \
		echo "Database $(DB_MARIADB_GATEWAY_DATABASE) already exists."; \
	fi

create-gateway-db-mariadb-user: check-db-mariadb create-gateway-db-mariadb
	@if ! $(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "SELECT User FROM mysql.user WHERE User='$(DB_MARIADB_GATEWAY_USER)';" | grep -q "$(DB_MARIADB_GATEWAY_USER)"; then \
		echo "Creating user $(DB_USER)..."; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "CREATE USER '$(DB_MARIADB_GATEWAY_USER)'@'%' IDENTIFIED BY '$(DB_MARIADB_GATEWAY_PASSWORD)';"; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "GRANT ALL PRIVILEGES ON $(DB_MARIADB_GATEWAY_DATABASE).* TO '$(DB_MARIADB_GATEWAY_USER)'@'%'; FLUSH PRIVILEGES;"; \
		echo "User $(DB_MARIADB_GATEWAY_USER) created."; \
	else \
		echo "User $(DB_MARIADB_GATEWAY_USER) already exists."; \
	fi

.PHONY: create-powerdns-db-mariadb
create-powerdns-db-mariadb: check-db-mariadb
	@if ! $(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "SHOW DATABASES LIKE '$(DB_MARIADB_PDNS_DATABASE)';" | grep -q "$(DB_MARIADB_PDNS_DATABASE)"; then \
		echo "Creating database $(DB_MARIADB_PDNS_DATABASE)..."; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "CREATE DATABASE $(DB_MARIADB_PDNS_DATABASE);"; \
		echo "Database $(DB_MARIADB_PDNS_DATABASE) created."; \
	else \
		echo "Database $(DB_MARIADB_PDNS_DATABASE) already exists."; \
	fi

.PHONY: create-powerdns-db-mariadb-user
create-powerdns-db-mariadb-user: check-db-mariadb create-powerdns-db-mariadb
	@if ! $(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "SELECT User FROM mysql.user WHERE User='$(DB_MARIADB_PDNS_USER)';" | grep -q "$(DB_MARIADB_PDNS_USER)"; then \
		echo "Creating user $(DB_USER)..."; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "CREATE USER '$(DB_MARIADB_PDNS_USER)'@'%' IDENTIFIED BY '$(DB_MARIADB_PDNS_PASSWORD)';"; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "GRANT ALL PRIVILEGES ON $(DB_MARIADB_PDNS_DATABASE).* TO '$(DB_MARIADB_PDNS_USER)'@'%'; FLUSH PRIVILEGES;"; \
		echo "User $(DB_MARIADB_PDNS_USER) created."; \
	else \
		echo "User $(DB_MARIADB_PDNS_USER) already exists."; \
	fi

check-db-postgres: up-postgres
	@echo "Checking if PostgreSQL is available..."
	@until $(POSTGRES_CMD) exec postgres pg_isready -U $(DB_POSTGRES_USER) > /dev/null 2>&1; do \
		echo "Waiting for PostgreSQL to accept connections..."; \
		sleep 1; \
	done
	@echo "PostgreSQL is up and accepting connections."

create-gateway-db-postgres: check-db-postgres
	@if ! $(POSTGRES_CMD) exec postgres psql -U $(DB_POSTGRES_USER) -tAc "SELECT 1 FROM pg_database WHERE datname='$(DB_POSTGRES_GATEWAY_DATABASE)';" | grep -q 1; then \
		echo "Creating database $(DB_POSTGRES_GATEWAY_DATABASE)..."; \
		$(POSTGRES_CMD) exec postgres psql -U $(DB_POSTGRES_USER) -c "CREATE DATABASE $(DB_POSTGRES_GATEWAY_DATABASE);"; \
		echo "Database $(DB_POSTGRES_GATEWAY_DATABASE) created."; \
	else \
		echo "Database $(DB_POSTGRES_GATEWAY_DATABASE) already exists."; \
	fi

create-gateway-db-postgres-user: check-db-postgres create-gateway-db-postgres
	@if ! $(POSTGRES_CMD) exec postgres psql -U $(DB_POSTGRES_USER) -tAc "SELECT 1 FROM pg_roles WHERE rolname='$(DB_POSTGRES_GATEWAY_USER)';" | grep -q 1; then \
		echo "Creating user $(DB_POSTGRES_GATEWAY_USER)..."; \
		$(POSTGRES_CMD) exec postgres psql -U $(DB_POSTGRES_USER) -c "CREATE USER $(DB_POSTGRES_GATEWAY_USER) WITH PASSWORD '$(DB_POSTGRES_GATEWAY_PASSWORD)';"; \
		$(POSTGRES_CMD) exec postgres psql -U $(DB_POSTGRES_USER) -c "GRANT ALL PRIVILEGES ON DATABASE $(DB_POSTGRES_GATEWAY_DATABASE) TO $(DB_POSTGRES_GATEWAY_USER);"; \
		echo "User $(DB_POSTGRES_GATEWAY_USER) created."; \
	else \
		echo "User $(DB_POSTGRES_GATEWAY_USER) already exists."; \
	fi

.PHONY: create-powerdns-db-postgres
create-powerdns-db-postgres: check-db-postgres
	@if ! $(POSTGRES_CMD) exec postgres psql -U $(DB_POSTGRES_USER) -tAc "SELECT 1 FROM pg_database WHERE datname='$(DB_POSTGRES_PDNS_DATABASE)';" | grep -q 1; then \
		echo "Creating database $(DB_POSTGRES_PDNS_DATABASE)..."; \
		$(POSTGRES_CMD) exec postgres psql -U $(DB_POSTGRES_USER) -c "CREATE DATABASE $(DB_POSTGRES_PDNS_DATABASE);"; \
		echo "Database $(DB_POSTGRES_PDNS_DATABASE) created."; \
	else \
		echo "Database $(DB_POSTGRES_PDNS_DATABASE) already exists."; \
	fi

.PHONY: create-powerdns-db-postgres-user
create-powerdns-db-postgres-user: check-db-postgres create-powerdns-db-postgres
	@if ! $(POSTGRES_CMD) exec postgres psql -U $(DB_POSTGRES_USER) -tAc "SELECT 1 FROM pg_roles WHERE rolname='$(DB_POSTGRES_PDNS_USER)';" | grep -q 1; then \
		echo "Creating user $(DB_POSTGRES_PDNS_USER)..."; \
		$(POSTGRES_CMD) exec postgres psql -U $(DB_POSTGRES_USER) -c "CREATE USER $(DB_POSTGRES_PDNS_USER) WITH PASSWORD '$(DB_POSTGRES_PDNS_PASSWORD)';"; \
		$(POSTGRES_CMD) exec postgres psql -U $(DB_POSTGRES_USER) -c "GRANT ALL PRIVILEGES ON DATABASE $(DB_POSTGRES_PDNS_DATABASE) TO $(DB_POSTGRES_PDNS_USER);"; \
		echo "User $(DB_POSTGRES_PDNS_USER) created."; \
	else \
		echo "User $(DB_POSTGRES_PDNS_USER) already exists."; \
	fi

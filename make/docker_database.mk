MARIADB_CMD=docker compose -f docker/docker-compose.mariadb.yml -p pdns-database

setup-mariadb-gateway: up-mariadb check-db create-gateway-db create-gateway-db-user

setup-mariadb-powerdns: up-mariadb check-db create-powerdns-db create-powerdns-db-user

setup-mariadb: setup-mariadb-gateway setup-mariadb-powerdns

check-db: up-mariadb
	@echo "Checking if MariaDB is available..."
	@until $(MARIADB_CMD) exec mariadb mariadb-admin -u $(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) ping > /dev/null 2>&1; do \
		echo "Waiting for MariaDB to accept connections..."; \
		sleep 1; \
	done
	@echo "MariaDB is up and accepting connections."

create-gateway-db: check-db
	@if ! $(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "SHOW DATABASES LIKE '$(DB_MARIADB_GATEWAY_DATABASE)';" | grep -q "$(DB_MARIADB_GATEWAY_DATABASE)"; then \
		echo "Creating database $(DB_MARIADB_GATEWAY_DATABASE)..."; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "CREATE DATABASE $(DB_MARIADB_GATEWAY_DATABASE);"; \
		echo "Database $(DB_MARIADB_GATEWAY_DATABASE) created."; \
	else \
		echo "Database $(DB_MARIADB_GATEWAY_DATABASE) already exists."; \
	fi

create-gateway-db-user: check-db create-gateway-db
	@if ! $(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "SELECT User FROM mysql.user WHERE User='$(DB_MARIADB_GATEWAY_USER)';" | grep -q "$(DB_MARIADB_GATEWAY_USER)"; then \
		echo "Creating user $(DB_USER)..."; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "CREATE USER '$(DB_MARIADB_GATEWAY_USER)'@'%' IDENTIFIED BY '$(DB_MARIADB_GATEWAY_PASSWORD)';"; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "GRANT ALL PRIVILEGES ON $(DB_MARIADB_GATEWAY_DATABASE).* TO '$(DB_MARIADB_GATEWAY_USER)'@'%'; FLUSH PRIVILEGES;"; \
		echo "User $(DB_MARIADB_GATEWAY_USER) created."; \
	else \
		echo "User $(DB_MARIADB_GATEWAY_USER) already exists."; \
	fi

.PHONY: create-powerdns-db
create-powerdns-db: check-db
	@if ! $(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "SHOW DATABASES LIKE '$(DB_MARIADB_PDNS_DATABASE)';" | grep -q "$(DB_MARIADB_PDNS_DATABASE)"; then \
		echo "Creating database $(DB_MARIADB_PDNS_DATABASE)..."; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "CREATE DATABASE $(DB_MARIADB_PDNS_DATABASE);"; \
		echo "Database $(DB_MARIADB_PDNS_DATABASE) created."; \
	else \
		echo "Database $(DB_MARIADB_PDNS_DATABASE) already exists."; \
	fi

.PHONY: create-powerdns-db-user
create-powerdns-db-user: check-db create-powerdns-db
	@if ! $(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "SELECT User FROM mysql.user WHERE User='$(DB_MARIADB_PDNS_USER)';" | grep -q "$(DB_MARIADB_PDNS_USER)"; then \
		echo "Creating user $(DB_USER)..."; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "CREATE USER '$(DB_MARIADB_PDNS_USER)'@'%' IDENTIFIED BY '$(DB_MARIADB_PDNS_PASSWORD)';"; \
		$(MARIADB_CMD) exec mariadb mariadb -u$(DB_MARIADB_ROOT_USER) -p$(DB_MARIADB_ROOT_PASSWORD) -e "GRANT ALL PRIVILEGES ON $(DB_MARIADB_PDNS_DATABASE).* TO '$(DB_MARIADB_PDNS_USER)'@'%'; FLUSH PRIVILEGES;"; \
		echo "User $(DB_MARIADB_PDNS_USER) created."; \
	else \
		echo "User $(DB_MARIADB_PDNS_USER) already exists."; \
	fi

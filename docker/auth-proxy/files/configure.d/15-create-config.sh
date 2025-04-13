#!/bin/bash

CONFIG_PATH=/etc/powerdns-auth-proxy

function create_jwt_config {
  if [ -z "${JWT_AUDIENCE}" ]; then
    echo "JWT_AUDIENCE is not set. Exiting."
    exit 1
  fi

  if [ -z "${JWT_ISSUER}" ]; then
    echo "JWT_ISSUER is not set. Exiting."
    exit 1
  fi

  if [ ! -z "${JWT_SECRET_KEY}" ] && [ ! -z "${JWT_PUBLIC_KEY}" ]; then
    JWT_MODE=0
  else
    JWT_MODE=1
  fi

  cat <<EOF >> ${CONFIG_PATH}/config.yaml
jwt:
  audience: ${JWT_AUDIENCE}
  issuer: ${JWT_ISSUER}
EOF

  if [ ${JWT_MODE} -eq 0 ]; then
cat <<EOF >> ${CONFIG_PATH}/config.yaml
  secret_key: ${JWT_SECRET_KEY}
  public_key: ${JWT_PUBLIC_KEY}
EOF
  else
cat <<EOF >> ${CONFIG_PATH}/config.yaml
  secret_key_path: ${JWT_SECRET_KEY_PATH:-/var/lib/powerdns-auth-proxy/data/jwt.privkey.pem}
  public_key_path: ${JWT_PUBLIC_KEY_PATH:-/var/lib/powerdns-auth-proxy/data/jwt.pubkey.pem}
EOF
  fi
}

function create_auth_config {
  echo "Configuring authentication type..."
  case "${AUTH_TYPE:-basic_auth}" in
    api_key)
      create_api_key_auth_config
      ;;
    basic_auth)
      create_basic_auth_config
      ;;
    jwt)
      create_jwt_auth_config
      ;;
    oauth)
      create_oauth_auth_config
      ;;
    *)
      echo "Unsupported auth type: ${AUTH_TYPE}. Exiting."
      exit 1
      ;;
  esac
}

function create_debug_config {
  if [ "${ENABLE_DEBUG}" == "true" ]; then
    DEBUG=true
  else
    DEBUG=false
  fi

  cat <<EOF > ${CONFIG_PATH}/config.yaml
debug: ${DEBUG}
EOF
}


function create_log_config {
  cat <<EOF >> ${CONFIG_PATH}/config.yaml
log_path: /dev/stdout
log_level: ${LOG_LEVEL:-info}
EOF
}

function create_api_key_auth_config {
  cat <<EOF >> ${CONFIG_PATH}/config.yaml
auth_type: api_key

EOF
}

function create_basic_auth_config {
  cat <<EOF >> ${CONFIG_PATH}/config.yaml
auth_type: basic_auth

EOF
}

function create_jwt_auth_config {
  cat <<EOF >> ${CONFIG_PATH}/config.yaml
auth_type: jwt

EOF
}

function create_oauth_auth_config {
  cat <<EOF >> ${CONFIG_PATH}/config.yaml
auth_type: oauth

EOF
}

function create_database_config {
  if [ -z "${DB_TYPE}" ]; then
    echo "DB_TYPE is not set. Exiting."
    exit 1
  fi

  case "${DB_TYPE}" in
    mariadb|mysql)
      create_mysql_config
      ;;
    postgres)
      create_postgres_config
      ;;
    sqlite)
      create_sqlite_config
      ;;
    *)
      echo "Unsupported database type: ${DB_TYPE}. Exiting."
      exit 1
      ;;
  esac
}

function create_mysql_config {
  if [ -z "${DB_MYSQL_HOST}" ]; then
    echo "DB_MYSQL_HOST is not set. Exiting."
    exit 1
  fi

  if [ -z "${DB_MYSQL_USER}" ]; then
    echo "DB_MYSQL_USER is not set. Exiting."
    exit 1
  fi

  if [ -z "${DB_MYSQL_PASSWORD}" ]; then
    echo "DB_MYSQL_PASSWORD is not set. Exiting."
    exit 1
  fi

  if [ -z "${DB_MYSQL_DATABASE}" ]; then
    echo "DB_MYSQL_DATABASE is not set. Exiting."
    exit 1
  fi

cat <<EOF >> ${CONFIG_PATH}/config.yaml
database: mysql

mysql:
  host: ${DB_MYSQL_HOST}
  port: ${DB_MYSQL_PORT:-3306}
  user: ${DB_MYSQL_USER}
  password: ${DB_MYSQL_PASSWORD}
  database: ${DB_MYSQL_DATABASE}
EOF
}

function create_postgres_config {
  if [ -z "${DB_POSTGRES_HOST}" ]; then
    echo "DB_POSTGRES_HOST is not set. Exiting."
    exit 1
  fi

  if [ -z "${DB_POSTGRES_USER}" ]; then
    echo "DB_POSTGRES_USER is not set. Exiting."
    exit 1
  fi

  if [ -z "${DB_POSTGRES_PASSWORD}" ]; then
    echo "DB_POSTGRES_PASSWORD is not set. Exiting."
    exit 1
  fi

  if [ -z "${DB_POSTGRES_DATABASE}" ]; then
    echo "DB_POSTGRES_DATABASE is not set. Exiting."
    exit 1
  fi

cat <<EOF >> ${CONFIG_PATH}/config.yaml
database: postgres

postgres:
  host: ${DB_POSTGRES_HOST}
  port: ${DB_POSTGRES_PORT:-5432}
  user: ${DB_POSTGRES_USER}
  password: ${DB_POSTGRES_PASSWORD}
  database: ${DB_POSTGRES_DATABASE}
EOF
}

function create_sqlite_config {
cat <<EOF >> ${CONFIG_PATH}/config.yaml
database: sqlite

sqlite:
  path: /var/lib/powerdns-auth-proxy/data/powerdns-auth-proxy.db
EOF
}

function create_pdns_config {
  if [ -z "${PDNS_HOST}" ]; then
    echo "PDNS_HOST is not set. Exiting."
    exit 1
  fi

  if [ -z "${PDNS_API_KEY}" ]; then
    echo "PDNS_API_KEY is not set. Exiting."
    exit 1
  fi

cat <<EOF >> ${CONFIG_PATH}/config.yaml

powerdns:
  host: ${PDNS_HOST}
  port: ${PDNS_PORT:-8081}
  ssl: ${PDNS_SSL:-false}
  api_key: ${PDNS_API_KEY}
EOF
}

mkdir -p ${CONFIG_PATH}
create_debug_config
create_log_config
create_auth_config
create_jwt_config
create_database_config
create_pdns_config

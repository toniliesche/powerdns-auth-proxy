# PowerDNS Authentication Proxy

## Introduction

This is a simple proxy that can be used to authenticate users against a PowerDNS server.

It allows to assign individual permissions to users on global or domain-based level.

## Building the Application

### building a local binary

The proxy is written in Golang and can be compiled by using the following command:

```bash
make build
```

Once the command has finished, you will find an executable file called "gateway" within your project directory.

To make this command work, you will need to have a running version of Go 1.22 or higher available, as well as the make
command.

### building a docker image

Although it is possible to run this application locally, it is highly recommended to run the service in Docker instead.

A Dockerfile based on Debian Bookworm and Go 1.22 is provided in the repository. To build the Docker image, just run the
following command:

```bash
make build-docker
```

To make this command work, you will need to have a running version of docker available, as well as the make command.

## Running the Application

### run locally

To run the application locally, you first have to create a yaml-formatted configuration file that can be used to setup
the application. The configuration file should be named `config.yaml` and could be placed anywhere on your local system.

Additionally you need to create ec signing key for your JWT based authentication which is mandatory for your admin
endpoints to work.

#### Create config.yaml

The following is an example of a configuration file that contains all available configuration keys. Read carefully and
extract the parts that match your desired setup for running the application.

```yaml
# true, false
debug: false

# path to directory where log files should be created
# e.g. /var/log/powerdns-auth-proxy
log_path: "<path to log file>"

# type of authentication to be used for PowerDNS proxy endpoints
# admin endpoints always use jwt auth
# valid values: api_key, basic_auth, jwt
auth_type: "<auth type>"


jwt:
  # audience string used for issued JWT tokens
  audience: "<audience>"
  # issuer string used for issued JWT tokens
  issuer: "<issuer>"
  # string value of the ec public key for JWT authentication
  # this MUST be combined with public_key
  secret_key: "<secret key>"
  # string value of the ec private key for JWT authentication
  # this MUST be combined with secret_key
  public_key: "<public key>"
  # path to public key file
  # this MUST be combined with public_key_path
  # e.g. /var/lib/powerdns-auth-proxy/data/jwt.privkey.pem
  secret_key_path: "<path to secret key file>"
  # path to public key file
  # this MUST be combined with secret_key_path
  # e.g. /var/lib/powerdns-auth-proxy/data/jwt.pubkey.pem
  public_key_path: "<path to public key file>"

# database type to be used
# valid values: mariadb, mysql, postgres, postgresql, sqlite
database: "<database type>"

# only necessary if database type is mariadb or mysql
mariadb:
  # hostname or ip of the mariadb server
  # e.g. 10.100.0.1, mariadb.example.com
  host: "<mariadb host>"
  # port of the mariadb server
  # e.g. 3306
  port: <mariadb port>
  # username for the mariadb database
  user: "<mariadb user>"
  # password for the mariadb database
  password: "<mariadb password>"
  # name of the mariadb database
  database: "<mariadb database>"

# only necessary if database type is postgres or postgresql
postgres:
  # hostname or ip of the postgres server
  # e.g. 10.100.0.1, postgres.example.com
  host: "<postgres host>"
  # port of the postgres server
  # e.g. 5432
  port: <postgres port>
  # username for the postgres database
  user: "<postgres user>"
  # password for the postgres database
  password: "<postgres password>"
  # name of the postgres database
  database: "<postgres database>"

# connection settings for powerdns authoritative server
powerdns:
  # hostname or ip of the powerdns server
  # e.g. 10.100.0.1, powerdns.example.com
  host: "<powerdns host>"
  # port of the powerdns server api
  # e.g. 8081
  port: <powerdns port>
  # setting if api is behind ssl
  ssl: <true/false>
  # api key for powerdns server
  # e.g. "changeme"
  api_key: "<powerdns api key>"
```

#### Create signing keys for jwt authentication

First we need to create a private key for signing the JWT tokens. This key should be kept secret and should not be
shared with anyone. The following command can be used to create a private key:

```bash
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out <path to private key file>
```

Afterwards we need to create a public key for verifying the JWT tokens. This key can be shared with anyone who wants to
verify the JWT tokens. The following command can be used to create a public key:

```bash
openssl ec -in <path to private key file> -pubout -out <path to public key file>
```

#### Run application

To run the application locally, you can use the following command:

```bash
./powerdns-auth-proxy --config <path to config> run
```

It might be necessary to create (or upgrade) the database schema before actually running the application, this can be
realized by the following command:

```bash
./powerdns-auth-proxy --config <path to config> db-migrate [--import-file <path to import file>]
```

### run with docker compose

To run the application in a Docker container, you need to create a docker-compose.yml file that can be used to setup the
application. The configuration file should be named `docker-compose.yml` and could be placed anywhere on your local
system.

#### Creating a docker-compose.yml

To run the application in a Docker container, you can use the following docker-compose.yml file as a template. To learn
more about the available configuration options, please refer to the section "environment variables".

```yaml
services:
  api:
    image: tliesche/powerdns-auth-proxy:latest
    environment:
      - AUTH_TYPE=jwt
      - DB_TYPE=sqlite
      - ENABLE_DEBUG=true
      - JWT_AUDIENCE=powerdns-auth-proxy
      - JWT_ISSUER=powerdns-auth-proxy
      - PDNS_API_KEY=changeme
      - PDNS_HOST=pdns-auth.example.com
      - PDNS_PORT=8081
    volumes:
      - proxy-data:/var/lib/powerdns-auth-proxy/data
    ports:
      - "8080:8080"

volumes:
  proxy-data:
    driver: local

networks:
  default:
    external: true
    name: pdns-network
```

#### environment variables

The following environment variables can be used to configure the application when running it in a Docker container. Read
carefully and extract the parts that match your desired setup for running the application.

| Variable             | Description                                                                                       | Required | Default | Values                                       |
|----------------------|---------------------------------------------------------------------------------------------------|----------|---------|----------------------------------------------|
| CA_CERTIFICATES_PATH | Path to the CA certificates directory                                                             | no       | `none`  | `string`                                     |
| DB_TYPE              | Type of the database to use                                                                       | yes      | `none`  | mariadb, mysql, postgres, postgresql, sqlite |
| DB_MYSQL_HOST        | Hostname of the MySQL/MariaDB server                                                              | yes¹     | `none`  | `string (fqdn or ip)`                        |
| DB_MYSQL_PORT        | Port of the MySQL/MariaDB server                                                                  | no       | 3306    | `integer`                                    |
| DB_MYSQL_NAME        | Name of the MySQL/MariaDB database                                                                | yes¹     | `none`  | `string`                                     |
| DB_MYSQL_USER        | Username for the MySQL/MariaDB database                                                           | yes¹     | `none`  | `string`                                     |
| DB_MYSQL_PASS        | Password for the MySQL/MariaDB database                                                           | yes¹     | `none`  | `string`                                     |
| DB_POSTGRES_HOST     | Hostname of the PostgreSQL server                                                                 | yes²     | `none`  | `string (fqdn or ip)`                        |
| DB_POSTGRES_PORT     | Port of the PostgreSQL server                                                                     | no       | 5432    | `integer`                                    |
| DB_POSTGRES_NAME     | Name of the PostgreSQL database                                                                   | yes²     | `none`  | `string`                                     |
| DB_POSTGRES_USER     | Username for the PostgreSQL database                                                              | yes²     | `none`  | `string`                                     |
| DB_POSTGRES_PASS     | Password for the PostgreSQL database                                                              | yes²     | `none`  | `string`                                     |
| ENABLE_DEBUG         | Enable debug mode                                                                                 | no       | false   | true, false                                  |
| JWT_AUDIENCE         | Audience string used for issued JWT tokens                                                        | yes      | `none`  | `string`                                     |
| JWT_ISSUER           | Issuer string used for issued JWT token                                                           | yes      | `none`  | `string`                                     |
| JWT_PUBLIC_KEY       | String value of the ec public key for JWT authentication, must be combined with `JWT_SECRET_KEY`  | no³      | `none`  | `string`                                     |
| JWT_SECRET_KEY       | String value of the ec private key for JWT authentication, must be combined with `JWT_PUBLIC_KEY` | no³      | `none`  | `string`                                     |
| JWT_PUBLIC_KEY_PATH  | Path to public key file, must be combined with `JWT_SECRET_KEY_PATH`                              | no³      | `none`  | `string`                                     |
| JWT_SECRET_KEY_PATH  | Path to secret key file, must be combined with `JWT_PUBLIC_KEY_PATH`                              | no³      | `none`  | `string`                                     |
| PDNS_API_KEY         | API key for connecting to the PowerDNS API                                                        | yes      | `none`  | `string`                                     |
| PDNS_HOST            | Hostname of the PowerDNS API                                                                      | yes      | `none`  | `string (fqdn or ip)`                        |
| PDNS_PORT            | Port of the PowerDNS API                                                                          | no       | 8081    | `integer`                                    |
| PDNS_SSL             | Use SSL for the PowerDNS API                                                                      | no       | false   | true, false                                  |

¹ Environment variables with prefix `DB_MYSQL_` are only mandatory when `DB_TYPE` is set to `mariadb`|`mysql`

² Environment variables with prefix `DB_POSTGRES_` are only mandatory when `DB_TYPE` is set to `postgres`|`postgresql`

³ JWT authentication setup:

- `SECRET_KEY` and `PUBLIC_KEY` are used when passing key values into the container directly
- `SECRET_KEY_PATH` and `PUBLIC_KEY_PATH` are used when passing pre-generated key files via mount into the container
- if none of these is defined, new keys will be generated on first container start and stored in
  `/var/lib/powerdns-auth-proxy/data/`

#### Run application

To run the application in a Docker container, you can use the following command:

```bash
docker-compose -f <path to your docker-compose.yaml> -p "<name of your docker compose stack>" up -d
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This project is not affiliated with, endorsed by, or sponsored by PowerDNS. It is an independent tool designed to act as
an authentication and authorization proxy for the PowerDNS HTTP API. All trademarks, service names, and product names
mentioned are the property of their respective owners. The use of "PowerDNS" is solely for descriptive purposes and does
not imply any association with or endorsement by PowerDNS. Users are responsible for ensuring compliance with PowerDNS
licensing and security requirements when using this tool.
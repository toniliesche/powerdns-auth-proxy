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

To run the application locally, you can use the following command:

```bash
./gateway --config <path to config> run
```

It might be necessary to create (or upgrade) the database schema before actually running the application, this can be
realized by the following command:

```bash
./gateway --config <path to config> db-migrate [--import-file <path to import file>]
```

### configuration file

```yaml
database: "<database type>" # mariadb, mysql, sqlite

```

### run with docker

### run with docker compose

### docker environment variables

| Variable             | Description                                | Required | Default | Values                 |
|----------------------|--------------------------------------------|----------|---------|------------------------|
| CA_CERTIFICATES_PATH | Path to the CA certificates directory      | no       | `none`  | `string`               |
| DB_TYPE              | Type of the database to use                | yes *    | `none`  | mariadb, mysql, sqlite |
| DB_MYSQL_HOST        | Hostname of the MySQL/MariaDB server       | yes *    | `none`  | `string (fqdn or ip)`  |
| DB_MYSQL_PORT        | Port of the MySQL/MariaDB server           | no       | 3306    | `integer`              |
| DB_MYSQL_NAME        | Name of the MySQL/MariaDB database         | yes *    | `none`  | `string`               |
| DB_MYSQL_USER        | Username for the MySQL/MariaDB database    | yes *    | `none`  | `string`               |
| DB_MYSQL_PASS        | Password for the MySQL/MariaDB database    | yes *    | `none`  | `string`               |
| PDNS_API_KEY         | API key for connecting to the PowerDNS API | yes      | `none`  | `string`               |
| PDNS_HOST            | Hostname of the PowerDNS API               | yes      | `none`  | `string (fqdn or ip)`  |
| PDNS_PORT            | Port of the PowerDNS API                   | no       | 8081    | `integer`              |
| PDNS_SSL             | Use SSL for the PowerDNS API               | no       | false   | true, false            |

* Only required if `DB_TYPE` is set to `mariadb` or `mysql`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This project is not affiliated with, endorsed by, or sponsored by PowerDNS. It is an independent tool designed to act as
an authentication and authorization proxy for the PowerDNS HTTP API. All trademarks, service names, and product names
mentioned are the property of their respective owners. The use of "PowerDNS" is solely for descriptive purposes and does
not imply any association with or endorsement by PowerDNS. Users are responsible for ensuring compliance with PowerDNS
licensing and security requirements when using this tool.
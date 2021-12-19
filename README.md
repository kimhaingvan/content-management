# content management

Headless cms for mobile and web apps

## Quick Start
- Install [Golang](https://golang.org/doc/install)
- Install [Docker](https://www.docker.com/)

### Set environment variables (example)

```bash
export APPLICATION_NAME=content-management
export CONSUL_IP=10.91.120.55
export CONSUL_PORT=8500
export CONSUL_ACL_TOKEN=7caf93ca-2112-2f84-3bc9-39e812983ed1
export LOGSTASH_IP=10.90.68.35
export LOGSTASH_PORT=30204
```

### Set consul config
- Setup consul config file at ${CONSUL_IP}:${CONSUL_PORT}/ui (Example: "config.json")
- { \
  &nbsp; &nbsp; "application_name": "otp-api",\
  &nbsp; &nbsp; "databases": {\
  &nbsp; &nbsp; &nbsp; &nbsp; "postgres_db": {\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "protocol":"localhost",\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "host":"localhost",\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;  "port":5432,\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "username":"postgres",\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "password":"postgres",\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "database":"postgres",\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "sslmode": "disable",\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "timeout": 15,\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "max_open_conns": 0,\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "max_conn_lifetime": 0,\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "google_auth_file": ""\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; }\
  &nbsp; &nbsp; &nbsp; &nbsp;},\
  &nbsp; &nbsp; "log": {\
  &nbsp; &nbsp; &nbsp; &nbsp; "level": "info"\
  &nbsp; &nbsp; },\
  &nbsp; &nbsp; "zipkin": {\
  &nbsp; &nbsp; &nbsp; &nbsp; "url": "http://10.90.68.35:30208"\
  &nbsp; &nbsp; },\
  }

### Build and start
```bash
$ docker-compose up -d
$ go run main.go
```

You should see the following message:

    HTTP server listening at :8080

To view CMS Administrator, open these URLs in browser:

- [http://localhost:8080/](http://localhost:8080/)

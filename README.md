# content management
main.go
Headless cms for mobile and web apps

## Quick Start
- Install [Docker](https://www.docker.com/)

### Set environment variables (example)

```bash
export APPLICATION_NAME=content-management
export CONSUL_IP=10.91.120.55
export CONSUL_PORT=8500
export CONSUL_ACL_TOKEN=7caf93ca-2112-2f84-3bc9-39e812983ed1
export LOGSTASH_IP=10.90.68.35
export LOGSTASH_PORT=30204
export ZIPKIN_URL=http://10.90.68.35:30208
export ELASTIC_APM_SERVER_URL=http://10.90.68.35:30207
export ELASTIC_APM_SERVICE_NAME=content-management
export ELASTIC_APM_ENVIRONMENT=development
export SERVER_PORT=8080
```

### Set consul config
- Setup consul config file at ${CONSUL_IP}:${CONSUL_PORT}/ui (Example: "config.json")
- { \
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
  &nbsp; &nbsp; "minio": {\
  &nbsp; &nbsp; &nbsp; &nbsp; "endpoint":"s3.cloud.cmctelecom.vn", \
  &nbsp; &nbsp; &nbsp; &nbsp; "access_key":"3S5MEJZLE2T2YCAG8ZWT", \
  &nbsp; &nbsp; &nbsp; &nbsp; "secret_access_key":"S0TNxHmo3GWMgEBeF2QAUMGtfRCVE7aVNCh3DXaL", \
  &nbsp; &nbsp; &nbsp; &nbsp; "bucket_name":"content-management-dev"\
  &nbsp; &nbsp; },\
  &nbsp; &nbsp; "log": {\
  &nbsp; &nbsp; &nbsp; &nbsp; "level": "info"\
  &nbsp; &nbsp; }
  
  }

### Build and start
```bash
$ docker-compose up -d
```

You should see the following message:

    HTTP server listening at :8080

To view CMS Administrator, open these URLs in browser:

- [http://localhost:8080/](http://localhost:8080/)

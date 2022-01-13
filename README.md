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
export SERVER_PORT=8011
```

### Set consul config
- Setup consul config file at ${CONSUL_IP}:${CONSUL_PORT}/ui (Example: "config.json")
- { \
  &nbsp; &nbsp; "databases": {\
  &nbsp; &nbsp; &nbsp; &nbsp; "postgres_db": {\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "host":"localhost",\
  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; "port":5432,\
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

### Build and run localhost
```bash
$ go run cmd/content-server/main.go
hoặc
$ ./recompile.bash
```

You should see the following message:

    HTTP server listening at :8011

To view CMS Administrator, open these URLs in browser:

- [http://localhost:8011/admin](http://localhost:8011/admin)

## Build with docker
- Để build với docker thì trước tiên phải cài đặt docker
- Thêm deamon docker
  {
  "registry-mirrors": [],
  "insecure-registries": [
  "10.91.120.43:8000",
  "repo.mafc.vn:8000"
  ],
  "debug": true,
  "experimental": false
  }
- Sử dụng tiếp các câu lệnh sau:
+ docker login 10.91.120.43:8000 sau đó sử dụng user/pass để login admin/admin123
+ docker repo.mafc.vn:8000 sau đó sử dụng user/pass để login
### Build image
docker build --rm -f Dockerfile -t {name image}:{version} .
+ Example : docker build --rm -f Dockerfile -t content:1.0.1 .
### Tag image
docker tag {name image}:{version}  {repo address}/{name image}:{version}
+ Example : docker tag content:1.0.1 10.91.120.43:8000/repository/mobile-project/content:1.0.1
### Push image
docker push {repo address}/{name image}:{version} đường dẫn vừa tag bên trên
+ Example : docker push 10.91.120.43:8000/repository/mobile-project/content:1.0.1
### Pull image
docker pull  {repo address}/{name image}:{version} đường dẫn vừa push bên trên
+ Example: docker pull 10.91.120.43:8000/repository/mobile-project/content:1.0.1
### Run with docker 
docker run -d -p 8011:8011 -e {env name}={ env value} {repo address}/{name image}:{version} 

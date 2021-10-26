# content management

Headless cms for mobile and web apps

## Quick Start
- Install [docker](https://www.docker.com/)
- Install [consul](https://www.consul.io/docs/install)

### Set environment variables
```bash
export CONSUL_CONFIG_KEY_VALUE=config/content-management/data
export CONSUL_IP_ADDRESS=127.0.0.1
export CONSUL_PORT=8500
```

### Build and start
```bash
$ docker-compose up -d
```

You should see the following message:

    HTTP server listening at :8080

To view CMS Administrator, open these URLs in browser:

- [http://localhost:8080/admin](http://localhost:8080/admin/)

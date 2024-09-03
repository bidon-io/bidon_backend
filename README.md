# BidOn
## Set up development environment
```shell
make local-init

docker compose up -d
```

### Manage migrations
```shell
go run ./cmd/bidon-migrate -help
```

### Start admin backend
```shell
go run ./cmd/bidon-admin
```

### Start sdkapi backend
```shell
go run ./cmd/bidon-sdkapi
```

### Run tests
```shell
make test
```

### Clean env
```shell
docker compose down --volumes --rmi local --remove-orphans || true
```

### Read from kafka
```shell
docker compose exec -it kafka kafka-console-consumer --bootstrap-server=localhost:9092 --topic=bidon-ad-events --from-beginning
```

## Start prod environment

Create personal account on https://maxmind.com.

Start Docker Compose:
```shell
MAXMIND_ACCOUNT_ID=<CHANGE_ME> \
MAXMIND_LICENSE_KEY=<CHANGE_ME> \
APP_SECRET=<CHANGE_ME> \
PG_PASSWORD=<CHANGE_ME> \
SNOWFLAKE_NODE_ID=<CHANGE_ME> \
docker compose -f docker-compose-prod.yml up -d
```

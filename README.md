### Setup development

```shell
gem install dip
cp -n bidon_api/.env.sample bidon_api/.env || true
cp -n bidon_back/.env.sample bidon_back/.env || true

cd bidon_back
dip provision
dip rails c
dip rails s
dip bash

cd bidon_api
dip provision
dip rails c
dip rails s
dip bash
```
#### Clean env
```shell
docker compose down --volumes --rmi local --remove-orphans || true
```

#### Read from kafka
```shell

docker compose exec -it kafka kafka-console-consumer --bootstrap-server=localhost:9092 --topic=bidon-ad-events --from-beginning
```

### Start prod environment
On `Mac M1` change `LD_PRELOAD: /usr/lib/aarch64-linux-gnu/libjemalloc.so` in `docker-compose-prod.yml`

Use the following command to generate `SECRET_KEY_BASE`:
```shell
docker compose -f docker-compose-prod.yml run --rm --no-deps bidon-backend rails secret
```
Create personal account on https://maxmind.com.

Start Docker Compose:
```shell
MAXMIND_ACCOUNT_ID=<CHANGE_ME> \
MAXMIND_LICENSE_KEY=<CHANGE_ME> \
SECRET_KEY_BASE=<CHANGE_ME> \
PG_PASSWORD=<CHANGE_ME> \
docker compose -f docker-compose-prod.yml up -d
```

### Setup development

```shell
gem install dip
cp -n bidon_api/.env.sample bidon_api/.env || true
cp -n bidon_back/.env.sample bidon_back/.env || true
docker compose up -d postgres redis rails_back rails_api --build
docker compose run --rm rails_api bin/setup
docker compose run --rm rails_back bin/setup
docker compose up -d nginx

cd bidon_back
dip rails c
dip rails s
dip bash

cd bidon_api
dip rails c
dip rails s
dip bash
```
#### Clean env
```shell
docker compose down --volumes --rmi local --remove-orphans || true
```

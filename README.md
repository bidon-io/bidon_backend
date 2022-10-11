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

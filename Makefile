REGISTRY ?= "ghcr.io/bidon-io"

.PHONY: test
test:
	docker compose run --rm go-test

docker-build-push-prod-admin:
	docker buildx build --platform linux/amd64 --provenance=false \
	--target bidon-admin --cache-to type=inline --cache-from $(REGISTRY)/bidon-admin \
	-t $(REGISTRY)/bidon-admin:$(TAG) -t $(REGISTRY)/bidon-admin:latest --push .

docker-build-push-prod-sdkapi:
	docker buildx build --platform linux/amd64 --provenance=false \
	--target bidon-sdkapi --cache-to type=inline --cache-from $(REGISTRY)/bidon-sdkapi \
	-t $(REGISTRY)/bidon-sdkapi:$(TAG) -t $(REGISTRY)/bidon-sdkapi:latest --push .

docker-build-push-prod-migrate:
	docker buildx build --platform linux/amd64 --provenance=false \
	--target bidon-migrate --cache-to type=inline --cache-from $(REGISTRY)/bidon-migrate \
	-t $(REGISTRY)/bidon-migrate:$(TAG) -t $(REGISTRY)/bidon-migrate:latest --push .

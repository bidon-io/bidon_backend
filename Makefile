REGISTRY ?= ghcr.io/bidon-io

.PHONY: test

init: update-submodules
	@cp -n .env.sample .env || true
	@cp -n .env.test.sample .env.test || true

install-deps:
	@brew ls --versions buf || brew install bufbuild/buf/buf@1.47.2
	@brew ls --versions pre-commit || brew install pre-commit
	@pre-commit install

local-init: init install-deps

update-submodules:
	git submodule update --remote --recursive

test:
	docker compose run --rm go-test

tags = -t $(REGISTRY)/$(TARGET):$(VERSION)
ifneq ($(TAG),)
	tags += -t $(REGISTRY)/$(TARGET):$(TAG)
endif

docker-build-push-prod:
	docker buildx build --platform linux/amd64 --provenance=false \
	--target $(TARGET) --cache-to type=inline --cache-from $(REGISTRY)/$(TARGET) \
	$(tags) --push .

docker-build-push-prod-admin: override TARGET=bidon-admin
docker-build-push-prod-admin: docker-build-push-prod

docker-build-push-prod-sdkapi: override TARGET=bidon-sdkapi
docker-build-push-prod-sdkapi: docker-build-push-prod

docker-build-push-prod-migrate: override TARGET=bidon-migrate
docker-build-push-prod-migrate: docker-build-push-prod

docker-build-push-prod-proxy: override TARGET=bidon-proxy
docker-build-push-prod-proxy: docker-build-push-prod

docker-build-push-prod-copilot:
	cd copilot && \
	uv run langgraph dockerfile Dockerfile && \
	docker buildx build --platform linux/amd64 --provenance=false \
		--cache-to type=inline --cache-from $(REGISTRY)/copilot \
		$(tags) --push .

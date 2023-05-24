REGISTRY_EXT = "ghcr.io/bidon-io"
REGISTRY_INT = "registry.appodeal.com/bidon"

docker-build-push-prod-api:
	cd bidon_api && \
	docker buildx build --platform linux/amd64,linux/arm64 --provenance=false --target=prod \
	--build-arg BUILDKIT_INLINE_CACHE=1 --cache-from $(REGISTRY_EXT)/bidon-api:latest \
	-t $(REGISTRY_INT)/api:$(TAG) -t $(REGISTRY_INT)/api:latest -t $(REGISTRY_EXT)/bidon-api:$(TAG) -t $(REGISTRY_EXT)/bidon-api:latest  --push .

docker-build-push-prod-back:
	cd bidon_back && \
	docker buildx build --platform linux/amd64,linux/arm64 --provenance=false --target=prod \
	--build-arg BUILDKIT_INLINE_CACHE=1 --cache-from $(REGISTRY_EXT)/bidon-back:latest \
	-t $(REGISTRY_INT)/back:$(TAG) -t $(REGISTRY_INT)/back:latest -t $(REGISTRY_EXT)/bidon-back:$(TAG) -t $(REGISTRY_EXT)/bidon-back:latest --push .
docker-build-push-prod-admin:
	docker buildx build --platform linux/amd64,linux/arm64 --provenance=false \
	--build-arg BIDON_SERVICE=bidon-admin --cache-to type=inline --cache-from $(REGISTRY_EXT)/bidon-admin \
	-t $(REGISTRY_INT)/admin:$(TAG) -t $(REGISTRY_INT)/admin:latest -t $(REGISTRY_EXT)/bidon-admin:$(TAG) -t $(REGISTRY_EXT)/bidon-admin:latest --push .

REGISTRY = "ghcr.io/bidon-io"

docker-build-push-prod-api:
	cd bidon_api && \
	docker buildx build --platform linux/amd64,linux/arm64 --provenance=false --target=prod \
	--build-arg BUILDKIT_INLINE_CACHE=1 --cache-from $(REGISTRY)/bidon-api:latest \
	-t $(REGISTRY)/bidon-api:$(TAG) -t $(REGISTRY)/bidon-api:latest  --push .

docker-build-push-prod-back:
	cd bidon_back && \
	docker buildx build --platform linux/amd64,linux/arm64 --provenance=false --target=prod \
	--build-arg BUILDKIT_INLINE_CACHE=1 --cache-from $(REGISTRY)/bidon-back:latest \
	-t $(REGISTRY)/bidon-back:$(TAG) -t $(REGISTRY)/bidon-back:latest --push .

docker-build-push-prod-admin:
	docker buildx build --platform linux/amd64,linux/arm64 --provenance=false \
	--target bidon-admin --cache-to type=inline --cache-from $(REGISTRY)/bidon-admin \
	-t $(REGISTRY)/bidon-admin:$(TAG) -t $(REGISTRY)/bidon-admin:latest --push .

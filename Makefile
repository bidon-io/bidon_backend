docker-build-prod-api:
	cd bidon_api && docker build --target=prod -t registry.appodeal.com/bidon/api:$(TAG) -t registry.appodeal.com/bidon/api:latest -t ghcr.io/bidon-io/bidon-api:latest -t ghcr.io/bidon-io/bidon-api:$(TAG) .

docker-push-prod-api:
	docker push registry.appodeal.com/bidon/api:$(TAG)
	docker push registry.appodeal.com/bidon/api:latest
	docker push ghcr.io/bidon-io/bidon-api:$(TAG)
	docker push ghcr.io/bidon-io/bidon-api:latest

docker-build-prod-back:
	cd bidon_back && docker build --target=prod -t registry.appodeal.com/bidon/back:$(TAG) -t registry.appodeal.com/bidon/back:latest -t ghcr.io/bidon-io/bidon-back:latest -t ghcr.io/bidon-io/bidon-back:$(TAG) .

docker-push-prod-back:
	docker push registry.appodeal.com/bidon/back:$(TAG)
	docker push registry.appodeal.com/bidon/back:latest
	docker push ghcr.io/bidon-io/bidon-back:$(TAG)
	docker push ghcr.io/bidon-io/bidon-back:latest

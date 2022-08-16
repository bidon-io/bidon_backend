docker-build-prod-api:
	cd bidon_api && docker build --target=prod  -t registry.appodeal.com/bidon/api:$(TAG) .

docker-push-prod-api:
	docker push registry.appodeal.com/bidon/api:$(TAG)

docker-build-prod-back:
	cd bidon_back && docker build --target=prod -t registry.appodeal.com/bidon/back:$(TAG) .

docker-push-prod-back:
	docker push registry.appodeal.com/bidon/back:$(TAG)

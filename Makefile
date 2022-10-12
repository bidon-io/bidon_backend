docker-build-prod-api:
	cd bidon_api && docker build --target=prod --build-arg AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) --build-arg AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY)  -t registry.appodeal.com/bidon/api:$(TAG)  .
	docker tag registry.appodeal.com/bidon/api:$(TAG) registry.appodeal.com/bidon/api:latest

docker-push-prod-api:
	docker push registry.appodeal.com/bidon/api:$(TAG)
	docker push registry.appodeal.com/bidon/api:latest

docker-build-prod-back:
	cd bidon_back && docker build --target=prod -t registry.appodeal.com/bidon/back:$(TAG) .
	docker tag registry.appodeal.com/bidon/back:$(TAG) registry.appodeal.com/bidon/back:latest

docker-push-prod-back:
	docker push registry.appodeal.com/bidon/back:$(TAG)
	docker push registry.appodeal.com/bidon/back:latest

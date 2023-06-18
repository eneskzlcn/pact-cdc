basket-service:
	docker-compose -f consumer/basket-service/docker-compose.yml up -d --wait \
	&& go run consumer/basket-service/main.go

pact-broker:
	docker-compose -f pact-broker/docker-compose.yml up -d --wait


product-service:
	docker-compose -f provider/product-service/docker-compose.yml up -d --wait \
	&& go run provider/product-service/main.go

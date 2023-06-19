basket-service:
	docker-compose -f basket-service/docker-compose.yml up -d --wait \
	&& go run basket-service/main.go

pact-broker:
	docker-compose -f pact-broker/docker-compose.yml up -d --wait

product-service:
	docker-compose -f product-service/docker-compose.yml up -d --wait \
	&& go run product-service/main.go

stock-service:
	docker-compose -f stock-service/docker-compose.yml up -d --wait \
	&& go run stock-service/main.go

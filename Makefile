basket-service:
	docker-compose -f consumer/basket-service/docker-compose.yml up -d --wait \
	&& go run consumer/basket-service/main.go
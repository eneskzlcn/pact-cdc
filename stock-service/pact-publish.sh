#!/bin/bash

set -x

VERSION=0.0.1 #like 1.0.0
BROKER_BASE_URL=http://localhost
TAG=dev
BRANCH=arf-1000-stock-service-pact

#curl -X PUT \
#    http://localhost/pacts/provider/ProductService/consumer/BasketService/version/${VERSION} \
#    -H "Content-Type: application/json" \
#    -d @/Users/eneskizilcin/Documents/go-projects/pact-cdc-test-/consumer/basket-service/app/product/pacts/basketservice-productservice.json

pact-broker publish \
./app/product/pacts/stockservice-productservice.json \
--consumer-app-version=${VERSION} \
--broker-base-url=${BROKER_BASE_URL} \
--tag=${TAG} \
--branch=${BRANCH}

#!/bin/bash

set -x

VERSION=1.0.0 #like 1.0.0

curl -X PUT \
    http://localhost:9292/pacts/provider/ProductService/consumer/BasketService/version/${VERSION} \
    -H "Content-Type: application/json" \
    -d @/Users/eneskizilcin/Documents/go-projects/pact-cdc-test-/consumer/basket-service/app/product/pacts/basketservice-productservice.json
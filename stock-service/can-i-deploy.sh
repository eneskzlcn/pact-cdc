#!/bin/bash

set -x

VERSION=0.0.1 #like 1.0.0
PACTICIPANT=StockService
#TAG=test
#ENVIRONMENT=test
BROKER_BASE_URL=http://localhost
#BRANCH=arf-687-pact-test-spike

pact-broker can-i-deploy \
--pacticipant=${PACTICIPANT} \
--version=${VERSION} \
--broker-base-url=${BROKER_BASE_URL}
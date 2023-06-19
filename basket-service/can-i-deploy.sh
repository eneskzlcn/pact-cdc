#!/bin/bash

set -x

VERSION=4.0.2 #like 1.0.0
PACTICIPANT=BasketService
#TAG=test
#ENVIRONMENT=test
BROKER_BASE_URL=http://localhost
#BRANCH=arf-687-pact-test-spike

pact-broker can-i-deploy \
--pacticipant=${PACTICIPANT} \
--version=${VERSION} \
--broker-base-url=${BROKER_BASE_URL}
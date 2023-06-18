#!/bin/bash

set -x

VERSION=4.0.2 #like 1.0.0
PACTICIPANT=BasketService
TAG=dev
BROKER_BASE_URL=http://localhost

pact-broker create-version-tag \
--pacticipant "${PACTICIPANT}" \
--version "${VERSION}" \
--tag "${TAG}" \
--broker-base-url "${BROKER_BASE_URL}"
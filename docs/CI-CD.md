
How Pact validation workflow automation works:
https://docs.pact.io/pact_nirvana/step_6


Can I Deploy

```bash

can-i-deploy --participant {{$CONSUMER/PROVIDER_NAME}} --version {{$CONSUMER_VERSION}} --to {{$ENVIRONMENT}}

```

```bash
can-i-deploy --pacticipant BasketService --version 1.0.1 --to prod
```


Pipeline Stages

- Run Consumer Test / Generate Pact File
- Publish Pact File With Consumer Version (WITH GIT_COMMIT_SHORT_HASH VERSION OR YOU CAN USE LATEST AS DEFAULT)
- Tag The Consumer With The Current Environment (if you are deploying to the dev, you need to tag the consumer as dev/development to provide environment specific validation between pacts. Ex: You deploy your consumer api to the dev environment, then you need to trigger the provider pipeline for development branch, or if you deploy to master, you need to trigger master branch etc.)
- Trigger Webhook To Run Provider Tests For Each Service The Current Service Consumes with specified tags and version
- Each trigger-able pipeline stage for that providers, should succeed to be able to deploy the consumer
- Each provider tests succeed/verified and the results published to the broker
- Use can-i-deploy stage to be sure all provider succeed you do not have any problem on cdc interactions.
- Now you can deploy the consumer.



NOTES

-we should not put /api/v1 parts to the urls incoming from environment variables, 
instead we should use actual urls like https://localhost:5001 from the environment variables, and the '/api/v1' parts on the path in the client that sends the request.
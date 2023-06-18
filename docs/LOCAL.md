
you need to install pact cli via gem package manager to run pact
consumer tests on your local;

You can download with using brew or sth. but this causes some problems on pact inner packages. The cleanest way seems gem.

```bash
  gem install pact
```

you need to install pact-provider-verifier executable via gem to run provider tests

```bash
    gem install pact-provider-verifier
```

you need to install pact_broker-client to create tags or webhooks, test webhooks etc. on
your local

```bash
    gem install pact_broker-client
```


Remember that a pacticipant version in the Pact Broker should map 1:1 to a commit in your repository. To facilitate this, the version number used to publish pacts and verification results should either be or contain the commit.
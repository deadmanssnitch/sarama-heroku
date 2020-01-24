# Sarama Heroku

[![GoDoc](https://godoc.org/github.com/deadmanssnitch/sarama-heroku?status.svg)](http://godoc.org/github.com/deadmanssnitch/sarama-heroku)

## Overview

sarama-heroku is a Go library that makes it easy to connect to [Apache Kafka on
Heroku](https://www.heroku.com/kafka).  We handle all the certificate
management and configuration so that you can start up your Kafka consumers and
producers with minimal effort.

## Installation

```console
go get -u github.com/deadmanssnitch/sarama-heroku
```

### Heroku

Make sure you have the Heroku CLI plugin installed:
```console
heroku plugins:install heroku-kafka
```

Next, you'll need to provision a new Kafka add-on or attach an existing one to
your app. To provision run:
```console
heroku addons:create heroku-kafka:basic-0 -a [app]
heroku kafka:wait -a [app]
```

## Consumers

Now you are ready to start using the library.

Create a cluster consumer config like the following:

```go
kfkCfg, err := heroku.NewConfig()

config := cluster.NewConfig()
config.ClientID = "app-name." + os.Getenv("DYNO")
config.Net.TLS.Enable = kfkCfg.TLS()
config.Net.TLS.Config = kfkCfg.TLSConfig()

groupID := kfkCfg.Prefix("group-id")
topics := []string{kfkCfg.Prefix("topic")}

consumer, err := cluster.NewClusterConsumer(groupID, topics, config)
```

:heavy_exclamation_mark: Multi-tenant plans require creating the consumer
groups before you can use them.
```console
heroku kafka:consumer-groups:create 'group-id' -a [app]
```

## Producers

Furthermore, a producer can be either Sync or Async. Read up on the differences
[here](https://godoc.org/github.com/Shopify/sarama).

Creating an async producer from a custom config:

```go
kfkCfg, err := heroku.NewConfig()

config := sarama.NewConfig()
config.Producer.Return.Errors = true
config.Producer.RequiredAcks = sarama.WaitForAll
config.Net.TLS.Enable = kfkCfg.TLS()
config.Net.TLS.Config = kfkCfg.TLSConfig()

producer, err := sarama.NewAsyncProducer(kfkCfg.Brokers(), config)
```

:heavy_exclamation_mark: Multi-tenant plans require adding the KAFKA_PREFIX
when sending messages. You should use heroku.AppendPrefixTo("topic") to ensure
it's set.

```go
producer <- &sarama.ProducerMessage{
  Topic: heroku.AppendPrefixTo("events"),
  Key:   sarama.StringEncoder(key),
  Value: []byte("Message"),
}
```
For more information about how to set up a config see the
[sarama documentation](http://godoc.org/github.com/Shopify/sarama#Config).

## Multiple Kafka Instances

Use `heroku.NewConfigWithName` to build a config for a named Kafka instance.

```go
kfkCfg, err := heroku.NewConfigWithName("ONYX")
```

## Environment

Sarama Heroku depends on the following environment variables that are set by
the Heroku Kafka add-on:

  - KAFKA\_CLIENT\_CERT
  - KAFKA\_CLIENT\_CERT\_KEY
  - KAFKA\_TRUSTED\_CERT
  - KAFKA\_PREFIX (only multi-tenant plans)
  - KAFKA\_URL

## Contributing

Thank you so much for your interest in contributing to this repository. We
appreciate you and the work you're doing on this SO much.

For details [see CONTRIBUTING.md](CONTRIBUTING.md)

## Thanks

This package was extracted from [Dead Man's Snitch](https://deadmanssnitch.com),
a dead simple monitoring service for Cron jobs and system liveness.

# Name TBD

### Overview

{name} is a Go library that makes connecting to an instance of Kafka on Heroku easy. 


## Installation

Add to your Go project  

```
import "github.com/collectiveidea/name"
```
Then, from your Go project directory:  

```
go get github.com/collectiveidea/name
```  
 
 Or, using gb:

```
gb vendor fetch github.com/collectiveidea/name
```
 

## Usage

First, export your Kafka credentials to Heroku.

Creating a cluster consumer using a custom config:

```
 config := cluster.NewConfig()
 config.Group.PartitionStrategy = cluster.StrategyRoundRobin
 config.Group.Return.Notifications = true
 config.ClientID = "consumer group"
 config.Consumer.Return.Errors = true
 
      
consumer, _ := goheroku.NewClusterConsumer("topic", "consumer group", config)
```
If you do not need to customize a config, you can easily create a cluster consumer from the default config using the follow code:

```
consumer, _ := goheroku.NewQuickClusterConsumer("topic", "consumer group")
```

Creating an async producer from a custom config:

```
config := sarama.NewConfig()
config.Producer.Return.Errors = true
config.Producer.RequiredAcks = sarama.WaitForAll
 
producer, _ := goheroku.NewAsyncProducer(config)
```

For more information about how to set up a config see the [sarama docs](http://godoc.org/github.com/Shopify/sarama#Config) or the [cluster docs](http://godoc.org/github.com/bsm/sarama-cluster#Config).



## Contributing
In the spirit of [free software](https://www.gnu.org/philosophy/free-sw.html), everyone is encouraged to help improve this project. Here are a few ways you can pitch in:


* [Report bugs](https://github.com/collectiveidea/name/issues).
* Fix bugs and submit [pull requests](https://github.com/collectiveidea/name/pulls).
* Write, clarify or fix documentation.
* Refactor code.

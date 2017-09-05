# Name TBD

### Overview

{name} is a Go library that makes connecting to an instance of Kafka on Heroku easy. We handle all the certificate management and configuration so that you can start up your Kafka consumers and producers with minimal effort.


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

First, export your Kafka credentials to Heroku. You'll want to make sure to include the following:

  - KAFKA\_CLIENT\_CERT
  - KAFKA\_CLIENT\_CERT\_KEY
  - KAFKA\_TRUSTED\_CERT
  - KAFKA\_PREFIX
  - KAFKA\_URL

#### Consumers
Now you are ready to start using the library. From here there are two options. You can either create a consumer using all the default configurations or create your own configuration.

If you do not need a custom cluster consumer, you can easily create a cluster consumer from the default config using the follow code:

```
consumer, err := goheroku.NewQuickClusterConsumer("topic", "consumer group")
```

Or you can create a cluster consumer using a custom config:

```
config := cluster.NewConfig()
config.Group.PartitionStrategy = cluster.StrategyRoundRobin
config.Group.Return.Notifications = true
config.ClientID = "consumer group"
config.Consumer.Return.Errors = true

consumer, err := goheroku.NewClusterConsumer("topic", "consumer group", config)
```
#### Producers
There are multiple options for producers. Similar to consumers, you are use the default config for a quick start option or you can specifiy your own config. Furthermore, a producer can be either Sync or Async. Read up on the differences [here](https://godoc.org/github.com/Shopify/sarama).

Creating an async producer from a custom config:

```
config := sarama.NewConfig()
config.Producer.Return.Errors = true
config.Producer.RequiredAcks = sarama.WaitForAll

producer, err := goheroku.NewAsyncProducer(config)
```

For more information about how to set up a config see the [sarama docs](http://godoc.org/github.com/Shopify/sarama#Config) or the [cluster docs](http://godoc.org/github.com/bsm/sarama-cluster#Config).



## Contributing
Thank you so much for your interest in contributing to this repository. We appreciate you and the work you're doing on this SO much.

**How often are we checking this repository for issues or PRs:**
If you post an issue or PR and do not hear back right away, do not worry! We aren't ignoring you. Expect to hear back within a couple of weeks. Typically, we check this repository once a month on the last Friday to review PRs, issues, and other maintenance duties.

**Submitting an Issue**

Before you contribute, please take a look at this [helpful article](https://opensource.guide/how-to-contribute/#how-to-submit-a-contribution) and follow the guidelines. It is important to supply enough context and information so that we can get to it as quickly as possible.

* Please include a screenshot
* Please include exact steps to replicate the bug or issue with as much information as possible
* Please include any "hunches" you may have about what the issue might be

**Submitting a PR:**

* Please [fork](https://help.github.com/articles/creating-a-pull-request-from-a-fork/) the respository, put your code change on a new branch and then create a pull request against master in the main repository.
* In the PR description, please include the following:
	- A link to the issue is there is one
	- A screenshot, gif, or detailed description of before and after functionality
	- Any applicable tests, if warranted
	- Instructions on how to QA if applicable
	
Thank you for contributing to open source software. Your work is helping us all make better software. Happy coding!

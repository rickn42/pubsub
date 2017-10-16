# Pubsub 

simple publish-subscribe system 

### how to use 

publisher
```go
var issue = "hello"

publisher := pubsub.NewPub()
publisher.Publish(issue, "example")
```

subscriber
```go
var subscriber pubsub.Subscriber
subscriber = pubsub.NewSubFunc(func(issue, v interface{}) {
	// do something
})
```

subscribe, unsubscribe
```go 
publisher.Subscribe(issue, subscriber)
//...

publisher.Unsubscribe(issue, subscriber)
```
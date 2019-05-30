package omni

import "time"

type Notification interface {
	GetMessageType() MessageType
	GetMessage() Message
}

type Initiator interface {
	GetSubscribers() <-chan Notification
}

type BloomFilter struct{}

type Once struct {
	// For every messages that are going to be sent, ensure that it is only send once for the specified time window.
	// E.g. If we want to send a recommendation daily, then the unique process id can be set to the unique message id + the timestamp start of the day.
	// The process id will be added to the bloom filter for every successful send.
	// When the server restart, we can also find the last id and resume from the delta position (given that the data is queried in a sorted key, e.g. user id in ascending order, then we can resume from a given id)
	bloomfilter BloomFilter
	lastID      string
	lastSent    time.Time
}

// Accepts a notification, checks the type and returns a factory content builder that is responsible for
// building the message.
type Delegator interface {
	FactoryContentBuilder(Notification)
}

type ContentBuilder interface {
	// Given the stream of subscriber, build the dynamic message and send the Message to another out channel.
	// If the message is static, we can just use the flyweight pattern and send a copy of the message
	// directly, skipping the ContentBuilder.
	Build(<-chan Notification) <-chan Message
}

type Dispatcher interface {
	Dispatch(Message) error
}

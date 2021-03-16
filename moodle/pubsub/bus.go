package pubsub

import "fmt"

// pubs is a mapping from a topic to the publisher output channel
var pubs map[Topic]chan interface{} = make(map[Topic]chan interface{})

// subs is a mapping from a topic to subscribed channels
var subs map[Topic][]chan interface{} = make(map[Topic][]chan interface{})

// Publish returns channel that can be written to
// Not concurrent safe
func Publish(topic Topic) (channel chan<- interface{}, err error) {
	if _, ok := subs[topic]; ok {
		err = fmt.Errorf("Publish topic already published %s", topic)
		return
	}

	subs[topic] = []chan interface{}{} // empty subcribers slice
	channel = make(chan interface{})

	return
}

// Close cleans up a topic, closing all subscribed channels
// Not concurrent safe
func Close(topic Topic) (err error) {
	if _, ok := subs[topic]; !ok {
		err = fmt.Errorf("Close: topic not published %s", topic)
		return
	}

	return
}

// Subscribe returns a read channel for the requested topic
// Not concurrent safe
func Subscribe(topic Topic) (channel <-chan interface{}, err error) {
	if _, ok := subs[topic]; !ok {
		err = fmt.Errorf("Subscribe: topic not published %s", topic)
		return
	}

	channel = make(chan interface{})
	return
}

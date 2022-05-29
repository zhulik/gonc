package eventbus

import (
	"sync/atomic"
)

type eventHandler func(event any) bool

type messageType int

const (
	messageTypeEvent messageType = iota
	messageTypeSubscription
	messageTypeUnsubscription
	messageTypeTopics
)

type message struct {
	messageType messageType

	topic string
	event any

	subscription *Subscription
	done         chan []string
}

type EventBus struct {
	inputChannel chan message

	subscriptions subscriptionStorage
	active        int32

	stopped chan bool
}

func New() EventBus {
	bus := EventBus{
		inputChannel:  make(chan message),
		stopped:       make(chan bool, 1),
		subscriptions: newSubscriptionStorage(),
		active:        1,
	}
	go bus.listen()

	return bus
}

type genericEventHandler[E any] func(event E) bool

func Subscribe[E any](b EventBus, topic string, handler genericEventHandler[E]) *Subscription {
	h := func(event any) bool {
		e, ok := event.(E)
		if !ok {
			panic(ok)
		}

		return handler(e)
	}

	return b.Subscribe(topic, h)
}

func (b *EventBus) Topics() []string {
	if !b.IsActive() {
		return []string{}
	}

	done := make(chan []string)
	msg := message{messageType: messageTypeTopics, done: done} // nolint: exhaustruct
	b.inputChannel <- msg

	return <-done
}

func (b EventBus) Publish(topic string, event any) {
	b.inputChannel <- message{messageType: messageTypeEvent, topic: topic, event: event} // nolint: exhaustruct
}

func (b *EventBus) Stop() {
	atomic.StoreInt32(&(b.active), 0)
	close(b.inputChannel)
	<-b.stopped
}

func (b *EventBus) IsActive() bool {
	return atomic.LoadInt32(&(b.active)) == 1
}

func (b EventBus) Subscribe(topic string, handler eventHandler) *Subscription {
	s := &Subscription{Handler: handler, Topic: topic, bus: b}

	done := make(chan []string)
	msg := message{messageType: messageTypeSubscription, subscription: s, done: done} // nolint: exhaustruct
	b.inputChannel <- msg
	<-done

	return s
}

func (b EventBus) Unsubscribe(s *Subscription) {
	done := make(chan []string)
	msg := message{messageType: messageTypeUnsubscription, subscription: s, done: done} // nolint: exhaustruct
	b.inputChannel <- msg
	<-done
}

func (b EventBus) listen() {
	for msg := range b.inputChannel {
		switch msg.messageType {
		case messageTypeEvent:
			b.broadcastEvent(msg)
		case messageTypeSubscription:
			b.subscriptions.add(msg.subscription)
			msg.done <- []string{}
		case messageTypeUnsubscription:
			b.subscriptions.remove(msg.subscription)
			msg.done <- []string{}
		case messageTypeTopics:
			msg.done <- b.subscriptions.topics()
		}
	}
	b.stopped <- true
}

func (b EventBus) broadcastEvent(msg message) {
	subs := b.subscriptions.get(msg.topic)
	if subs == nil {
		return
	}

	for subscription := range subs {
		if subscription.Handler(msg.event) {
			b.subscriptions.remove(subscription)
		}
	}
}

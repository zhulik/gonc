package eventbus

import (
	"github.com/zhulik/gonc/flag"
	"github.com/zhulik/gonc/notification"
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
	queue chan message

	subscriptions subscriptionStorage
	stoppedFlag   flag.Flag

	stopped notification.Notification
}

func New() EventBus {
	bus := EventBus{
		queue:         make(chan message),
		stopped:       notification.New(),
		subscriptions: newSubscriptionStorage(),
		stoppedFlag:   flag.New(),
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
	b.queue <- msg

	return <-done
}

func (b EventBus) Publish(topic string, event any) {
	b.queue <- message{messageType: messageTypeEvent, topic: topic, event: event} // nolint: exhaustruct
}

func (b *EventBus) Stop() {
	if !b.IsActive() {
		return
	}
	b.stoppedFlag.Raise()
	close(b.queue)
}

func (b EventBus) Wait() {
	b.stopped.Wait()
}

func (b *EventBus) StopWait() {
	b.Stop()
	b.Wait()
}

func (b *EventBus) IsActive() bool {
	return !b.stoppedFlag.IsRaised()
}

func (b EventBus) Subscribe(topic string, handler eventHandler) *Subscription {
	s := &Subscription{Handler: handler, Topic: topic, bus: b}

	done := make(chan []string)
	msg := message{messageType: messageTypeSubscription, subscription: s, done: done} // nolint: exhaustruct
	b.queue <- msg
	<-done

	return s
}

func (b EventBus) Unsubscribe(s *Subscription) {
	done := make(chan []string)
	msg := message{messageType: messageTypeUnsubscription, subscription: s, done: done} // nolint: exhaustruct
	b.queue <- msg
	<-done
}

func (b EventBus) listen() {
	for msg := range b.queue {
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
	b.stopped.Signal()
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

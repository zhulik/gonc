package eventbus

import "sort"

type handlerStorage map[*Subscription]*Subscription

type subscriptionStorage struct {
	storage map[string]handlerStorage
}

func newSubscriptionStorage() subscriptionStorage {
	return subscriptionStorage{
		storage: map[string]handlerStorage{},
	}
}

func (s subscriptionStorage) topics() []string {
	i := 0
	topics := make([]string, len(s.storage))
	for topic := range s.storage {
		topics[i] = topic
		i++
	}
	sort.Strings(topics)

	return topics
}

func (s *subscriptionStorage) add(sub *Subscription) {
	if _, ok := s.storage[sub.Topic]; !ok {
		s.storage[sub.Topic] = handlerStorage{}
	}
	s.storage[sub.Topic][sub] = sub
}

func (s subscriptionStorage) get(topic string) handlerStorage {
	return s.storage[topic]
}

func (s *subscriptionStorage) remove(sub *Subscription) {
	if _, ok := s.storage[sub.Topic]; !ok {
		return
	}
	if _, ok := s.storage[sub.Topic][sub]; !ok {
		return
	}
	delete(s.storage[sub.Topic], sub)
	if len(s.storage[sub.Topic]) == 0 {
		delete(s.storage, sub.Topic)
	}
}

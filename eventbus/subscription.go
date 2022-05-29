package eventbus

type Subscription struct {
	Topic   string
	Handler eventHandler
	bus     EventBus
}

func (s *Subscription) Unsubscribe() {
	s.bus.Unsubscribe(s)
}

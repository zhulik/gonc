package notification

import "sync"

type Notification interface {
	Signal()
	Wait()
}

type notifikation struct {
	wg sync.WaitGroup
}

func New() Notification {
	n := notifikation{
		wg: sync.WaitGroup{},
	}
	n.wg.Add(1)
	return &n
}

func (m *notifikation) Signal() {
	m.wg.Done()
}

func (m *notifikation) Wait() {
	m.wg.Wait()
}

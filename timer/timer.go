package timer

import (
	"time"
)

type Timer interface {
	Stop()
}

type timer struct {
	timer *time.Timer
}

func (t *timer) Stop() {
	t.timer.Stop()
}

func Every(d time.Duration, handler func() bool) Timer {
	t := time.NewTimer(d)
	go func() {
		for range t.C {
			if handler() {
				return
			}
		}
	}()
	return &timer{
		timer: t,
	}
}

func In(d time.Duration, handler func()) Timer {
	t := time.NewTimer(d)
	go func() {
		<-t.C
		handler()
	}()
	return &timer{
		timer: t,
	}
}

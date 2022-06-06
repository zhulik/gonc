package fspoller

import (
	"os"
	"sync"
	"time"

	"github.com/zhulik/gonc/timer"
)

type FSPoller interface {
	Stop()
	Watch(string, UpdateHadler) error
}

type UpdateHadler func(path string) bool

type watch struct {
	path         string
	handlers     []UpdateHadler
	lastModified time.Time
}

type fspoller struct {
	timer       timer.Timer
	watchesLock sync.Mutex
	watches     map[string]watch
}

func (p *fspoller) Stop() {
	p.timer.Stop()
}

func New(d time.Duration) FSPoller {
	poller := &fspoller{
		watchesLock: sync.Mutex{},
		watches:     make(map[string]watch),
	}
	handler := func() bool {
		poller.watchesLock.Lock()
		defer poller.watchesLock.Unlock()
		for path, watch := range poller.watches {
			info, err := os.Stat(path)
			if err != nil {
				panic(err)
			}
			if watch.lastModified == info.ModTime() {
				continue
			}
			watch.lastModified = info.ModTime()
			poller.watches[path] = watch

			wg := sync.WaitGroup{}
			wg.Add(len(watch.handlers))

			for _, h := range watch.handlers {
				go func(p string, handler UpdateHadler) {
					if handler(p) {
						// TODO: unsubscribe
					}
					wg.Done()
				}(path, h)
			}
			wg.Wait()

		}
		return false
	}
	timer := timer.Every(d, handler)
	poller.timer = timer
	return poller
}

func (p *fspoller) Watch(path string, handler UpdateHadler) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	info.ModTime()
	p.watchesLock.Lock()
	defer p.watchesLock.Unlock()
	// TODO: normalize path
	w, ok := p.watches[path]
	if !ok {
		w = watch{path: path, lastModified: info.ModTime(), handlers: []UpdateHadler{}}
		p.watches[path] = w
	}

	w.handlers = append(p.watches[path].handlers, handler)
	p.watches[path] = w
	return nil
}

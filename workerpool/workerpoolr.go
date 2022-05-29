package workerpool

import (
	"sync"
	"sync/atomic"

	"github.com/zhulik/gonc/future"
	"github.com/zhulik/gonc/notification"
)

type TaskR[T any] func() (T, error)

type messageR[T any] struct {
	task   TaskR[T]
	future future.FR[T]
}

type WorkerPoolR[T any] struct {
	queue   chan messageR[T]
	stopped notification.Notification
	size    int
	active  int32
}

func NewR[T any](size, queueSize int) WorkerPoolR[T] {
	pool := WorkerPoolR[T]{
		queue:   make(chan messageR[T], queueSize),
		stopped: notification.New(),
		size:    size,
		active:  1,
	}
	go pool.run()
	return pool
}

func (p WorkerPoolR[T]) Go(t TaskR[T]) future.FR[T] {
	f := future.NewR[T]()
	m := messageR[T]{task: t, future: f}
	p.queue <- m

	return f
}

func (p *WorkerPoolR[T]) Stop() {
	if !p.IsActive() {
		return
	}
	atomic.StoreInt32(&(p.active), 0)
	close(p.queue)
}

func (p *WorkerPoolR[T]) Wait() {
	p.stopped.Wait()
}

func (p *WorkerPoolR[T]) StopWait() {
	p.Stop()
	p.Wait()
}

func (p *WorkerPoolR[T]) IsActive() bool {
	return atomic.LoadInt32(&(p.active)) == 1
}

func (p *WorkerPoolR[T]) run() {
	wg := sync.WaitGroup{}
	wg.Add(p.size)
	for i := 0; i < p.size; i++ {
		go workerR(i, p.queue, &wg)
	}
	wg.Wait()
	p.stopped.Signal()
}

func workerR[T any](id int, queue <-chan messageR[T], wg *sync.WaitGroup) {
	for msg := range queue {
		msg.future.Resolve(msg.task())
	}
	wg.Done()
}

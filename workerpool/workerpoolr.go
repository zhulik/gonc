package workerpool

import (
	"github.com/zhulik/gonc/future"
)

type TaskR[T any] func() (T, error)

type WorkerPoolR[T any] struct {
	pool WorkerPool
}

func NewR[T any](size, queueSize int) WorkerPoolR[T] {
	return WorkerPoolR[T]{
		pool: New(size, queueSize),
	}
}

func (p WorkerPoolR[T]) Go(t TaskR[T]) future.FR[T] {
	f := future.NewR[T]()
	h := func() {
		f.Resolve(t())
	}
	p.pool.Go(h)
	return f
}

func (p *WorkerPoolR[T]) Stop() {
	p.pool.Stop()
}

func (p *WorkerPoolR[T]) Wait() {
	p.pool.Wait()
}

func (p *WorkerPoolR[T]) StopWait() {
	p.Stop()
	p.Wait()
}

func (p *WorkerPoolR[T]) IsActive() bool {
	return p.pool.IsActive()
}

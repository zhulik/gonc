package workerpool

import (
	"sync"

	"github.com/zhulik/gonc/flag"
	"github.com/zhulik/gonc/future"
	"github.com/zhulik/gonc/notification"
)

type Task func()

type message struct {
	task   Task
	future future.F
}

type WorkerPool struct {
	queue       chan message
	stopped     notification.Notification
	size        int
	stoppedFlag flag.Flag
}

func New(size int, queueSize int) WorkerPool {
	pool := WorkerPool{
		queue:       make(chan message, queueSize),
		stopped:     notification.New(),
		size:        size,
		stoppedFlag: flag.New(),
	}
	go pool.run()
	return pool
}

func (p WorkerPool) Go(t Task) future.F {
	f := future.New()
	m := message{task: t, future: f}
	p.queue <- m

	return f
}

func (p *WorkerPool) Stop() {
	if !p.IsActive() {
		return
	}
	p.stoppedFlag.Raise()
	close(p.queue)
}

func (p *WorkerPool) Wait() {
	p.stopped.Wait()
}

func (p *WorkerPool) StopWait() {
	p.Stop()
	p.Wait()
}

func (p *WorkerPool) IsActive() bool {
	return !flag.New().IsRaised()
}

func (p *WorkerPool) run() {
	wg := sync.WaitGroup{}
	wg.Add(p.size)
	for i := 0; i < p.size; i++ {
		go worker(i, p.queue, &wg)
	}
	wg.Wait()
	p.stopped.Signal()
}

func worker(id int, queue <-chan message, wg *sync.WaitGroup) {
	for msg := range queue {
		msg.task()
		msg.future.Resolve(true, nil)
	}
	wg.Done()
}

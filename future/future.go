package future

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Future[T any] interface {
	Resolve(T, error)
	IsResolved() bool
	Wait() (T, error)
}

type result[T any] struct {
	result T
	err    error
}

type future[T any] struct {
	done   *sync.WaitGroup
	result unsafe.Pointer
}

type handler func()
type handlerR[T any] func() (T, error)

type F Future[bool]
type FR[T any] Future[T]

func New() F {
	wg := sync.WaitGroup{}
	wg.Add(1)
	return &future[bool]{
		done: &wg,
	}
}

func NewR[T any]() FR[T] {
	wg := sync.WaitGroup{}
	wg.Add(1)
	return &future[T]{
		done: &wg,
	}
}

func Go(h handler) F {
	f := New()
	go func() {
		h()
		f.Resolve(true, nil)
	}()
	return f
}

func GoR[T any](h handlerR[T]) FR[T] {
	f := NewR[T]()
	go func() {
		f.Resolve(h())
	}()
	return f
}

func (f *future[T]) Resolve(res T, err error) {
	if f.IsResolved() {
		return
	}
	r := result[T]{result: res, err: err}
	atomic.StorePointer(&f.result, unsafe.Pointer(&r))
	f.done.Done()
}

func (f *future[T]) Wait() (T, error) {
	f.done.Wait()
	return ((*result[T])(f.result)).result, ((*result[T])(f.result).err)
}

func (f *future[T]) IsResolved() bool {
	return atomic.LoadPointer(&f.result) != nil
}

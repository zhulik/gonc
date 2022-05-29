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

type handler[T any] func() (T, error)

func New[T any]() Future[T] {
	wg := sync.WaitGroup{}
	wg.Add(1)
	return &future[T]{
		done: &wg,
	}
}

func F[T any](h handler[T]) Future[T] {
	f := New[T]()
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

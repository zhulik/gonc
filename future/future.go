package future

import (
	"sync/atomic"
	"unsafe"

	"github.com/zhulik/gonc/notification"
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
	done   notification.Notification
	result unsafe.Pointer
}

type handler func()
type handlerR[T any] func() (T, error)

type F Future[bool]
type FR[T any] Future[T]

func New() F {
	return &future[bool]{
		done: notification.New(),
	}
}

func NewR[T any]() FR[T] {
	return &future[T]{
		done: notification.New(),
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
	f.done.Signal()
}

func (f *future[T]) Wait() (T, error) {
	f.done.Wait()
	return ((*result[T])(f.result)).result, ((*result[T])(f.result).err)
}

func (f *future[T]) IsResolved() bool {
	return atomic.LoadPointer(&f.result) != nil
}

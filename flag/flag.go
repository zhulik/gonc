package flag

import "sync/atomic"

type Flag interface {
	IsRaised() bool
	Raise()
}

func New() Flag {
	return &flag{}
}

type flag struct {
	raised int32
}

func (f *flag) IsRaised() bool {
	return atomic.LoadInt32(&(f.raised)) == 1
}

func (f *flag) Raise() {
	atomic.StoreInt32(&(f.raised), 1)
}

func (f *flag) Lower() {
	atomic.StoreInt32(&(f.raised), 0)
}

package workerpool_test

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhulik/gonc/workerpool"
)

func TestWorkerpoolNoWait(t *testing.T) {
	pool := workerpool.New(10, 1)

	var total int32 = 0

	for i := 0; i < 10_000; i++ {
		pool.Go(func() {
			atomic.AddInt32(&total, 1)
			pool.IsActive()
		})
	}
	pool.StopWait()
	assert.Equal(t, int32(10_000), total)
}

func TestWorkerpoolWait(t *testing.T) {
	pool := workerpool.New(10, 1)

	total := 0 // No race conditionios because tasks are executed virtually in sequence

	for i := 0; i < 10_000; i++ {
		pool.Go(func() {
			total += 1
			pool.IsActive()
		}).Wait()
	}
	pool.StopWait()
	assert.Equal(t, 10_000, total)
}

package workerpool_test

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhulik/gonc/workerpool"
)

func TestWorkerPoolRNoWait(t *testing.T) {
	pool := workerpool.NewR[int](10, 1)

	var total int32 = 0

	for i := 0; i < 10_000; i++ {
		pool.Go(func() (int, error) {
			atomic.AddInt32(&total, 1)
			pool.IsActive()
			return 0, nil
		})
	}
	pool.StopWait()
	assert.Equal(t, int32(10_000), total)
}

func TestWorkerPoolRWait(t *testing.T) {
	pool := workerpool.NewR[int](10, 1)

	total := 0
	for i := 0; i < 10_000; i++ {
		n, _ := pool.Go(func() (int, error) {
			return 1, nil
		}).Wait()
		total += n
	}
	pool.StopWait()
	assert.Equal(t, 10_000, total)
}

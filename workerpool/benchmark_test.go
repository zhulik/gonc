package workerpool_test

import (
	"sync/atomic"
	"testing"

	"github.com/zhulik/gonc/workerpool"
)

func BenchmarkWorkerPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pool := workerpool.New(10, 1)

		var total int32 = 0

		for i := 0; i < 10_000; i++ {
			pool.Go(func() {
				atomic.AddInt32(&total, 1)
				pool.IsActive()
			})
		}
		pool.StopWait()
		if total != 10_000 {
			b.Fail()
		}
	}
}

func BenchmarkWorkerPoolR(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
		if total != 10_000 {
			b.Fail()
		}
	}
}

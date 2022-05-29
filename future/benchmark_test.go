package future_test

import (
	"testing"

	"github.com/zhulik/gonc/future"
)

func BenchmarkFuture(b *testing.B) {
	for i := 0; i < b.N; i++ {
		futures := make([]future.Future[bool], 1000)
		for i := 0; i < 1000; i++ {
			futures[i] = future.Go(func() {
			})
			futures[i].IsResolved()
		}
		for i := 0; i < 1000; i++ {
			futures[i].Wait()
		}
	}
}

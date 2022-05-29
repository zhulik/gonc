package future_test

import (
	"testing"

	"github.com/zhulik/gonc/future"
)

func BenchmarkFuture(b *testing.B) {
	for i := 0; i < b.N; i++ {
		futures := make([]future.Future[int], 1000)
		for i := 0; i < 1000; i++ {
			futures[i] = future.F(func() (int, error) {
				return 1, nil
			})
			futures[i].IsResolved()
		}
		for i := 0; i < 1000; i++ {
			futures[i].Wait()
		}
	}
}

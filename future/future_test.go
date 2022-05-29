package future_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhulik/gonc/future"
)

func TestFuture(t *testing.T) {
	futures := make([]future.Future[int], 1000)
	for i := 0; i < 1000; i++ {
		futures[i] = future.F(func() (int, error) {
			return 1, nil
		})
		futures[i].IsResolved()
	}
	wg := sync.WaitGroup{}
	wg.Add(1000 * 10)
	for i := 0; i < 1000; i++ {
		for j := 0; j < 10; j++ {
			go func(n int) {
				r, _ := futures[n].Wait()
				assert.Equal(t, 1, r)
				wg.Done()
			}(i)
		}
	}
	wg.Wait()

	for i := 0; i < 1000; i++ {
		assert.True(t, futures[i].IsResolved())
	}
}

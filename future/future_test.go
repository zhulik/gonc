package future_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhulik/gonc/future"
)

func TestFuture(t *testing.T) {
	futures := make([]future.Future[bool], 1000)
	for i := 0; i < 1000; i++ {
		futures[i] = future.Go(func() {

		})
		futures[i].IsResolved()
	}
	wg := sync.WaitGroup{}
	wg.Add(1000 * 10)
	for i := 0; i < 1000; i++ {
		for j := 0; j < 10; j++ {
			go func(n int) {
				futures[n].Wait()
				wg.Done()
			}(i)
		}
	}
	wg.Wait()

	for i := 0; i < 1000; i++ {
		assert.True(t, futures[i].IsResolved())
	}
}

package eventbus_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhulik/gonc/eventbus"
)

func TestNew(t *testing.T) {
	t.Parallel()
	bus := eventbus.New()

	assert.True(t, bus.IsActive())
	assert.Empty(t, bus.Topics())
	bus.StopWait()
	assert.False(t, bus.IsActive())
}

func TestEventDelivery(t *testing.T) {
	t.Run("no listeners", func(t *testing.T) {
		t.Parallel()
		bus := eventbus.New()
		// It does not block, nothing gets delivered
		bus.Publish("topic1", "message")
		bus.Publish("topic2", "message")
		bus.Publish("topic3", "message")
		bus.Publish("topic4", "message")
		assert.Empty(t, bus.Topics())
		bus.StopWait()
		assert.False(t, bus.IsActive())
	})

	t.Run("few listeners, no publications", func(t *testing.T) {
		t.Parallel()
		bus := eventbus.New()

		sub11 := eventbus.Subscribe(bus, "topic1", func(event string) bool { return false })
		sub12 := eventbus.Subscribe(bus, "topic1", func(event string) bool { return false })
		sub21 := eventbus.Subscribe(bus, "topic2", func(event string) bool { return false })
		sub22 := eventbus.Subscribe(bus, "topic2", func(event string) bool { return false })

		assert.Equal(t, []string{"topic1", "topic2"}, bus.Topics())

		sub11.Unsubscribe()
		sub12.Unsubscribe()

		sub21.Unsubscribe()
		sub22.Unsubscribe()

		assert.Empty(t, bus.Topics())

		bus.StopWait()
		assert.False(t, bus.IsActive())
	})

	t.Run("few listeners, few publications", func(t *testing.T) {
		t.Parallel()
		bus := eventbus.New()

		var topic1Count int32
		var topic2Count int32

		eventbus.Subscribe(bus, "topic1", func(event string) bool {
			assert.Equal(t, "topic1", event)
			atomic.AddInt32(&topic1Count, 1)

			return true
		})
		eventbus.Subscribe(bus, "topic1", func(event string) bool {
			assert.Equal(t, "topic1", event)
			atomic.AddInt32(&topic1Count, 1)

			return true
		})
		eventbus.Subscribe(bus, "topic2", func(event string) bool {
			assert.Equal(t, "topic2", event)
			atomic.AddInt32(&topic2Count, 1)

			return true
		})
		eventbus.Subscribe(bus, "topic2", func(event string) bool {
			assert.Equal(t, "topic2", event)
			atomic.AddInt32(&topic2Count, 1)

			return true
		})

		assert.Equal(t, []string{"topic1", "topic2"}, bus.Topics())

		bus.Publish("topic1", "topic1")
		bus.Publish("topic2", "topic2")

		bus.StopWait()

		assert.Equal(t, int32(2), topic1Count)
		assert.Equal(t, int32(2), topic2Count)

		assert.Empty(t, bus.Topics())
		assert.False(t, bus.IsActive())
	})
}

func TestConcurrency(t *testing.T) {
	t.Parallel()
	bus := eventbus.New()

	m := sync.Mutex{}
	counts := map[string]int64{"topic1": 0, "topic2": 0, "total": 0}

	for i := 0; i < 4; i++ {
		var topic string
		if i < 2 {
			topic = "topic1"
		} else {
			topic = "topic2"
		}
		eventbus.Subscribe(bus, topic, func(event any) bool {
			m.Lock()
			defer m.Unlock()
			counts[topic] += int64(event.(int))   // nolint: forcetypeassert
			counts["total"] += int64(event.(int)) // nolint: forcetypeassert

			return false
		})
	}

	assert.Equal(t, []string{"topic1", "topic2"}, bus.Topics())

	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				var topic string

				if j < 50 {
					topic = "topic1"
				} else {
					topic = "topic2"
				}
				bus.Publish(topic, j)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	bus.StopWait()

	assert.Empty(t, bus.Topics())

	assert.Equal(t, int64(24500), counts["topic1"])
	assert.Equal(t, int64(74500), counts["topic2"])
	assert.Equal(t, int64(99000), counts["total"])

	assert.False(t, bus.IsActive())
}

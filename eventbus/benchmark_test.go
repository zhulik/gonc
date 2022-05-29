package eventbus_test

import (
	"math/rand"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/zhulik/gonc/eventbus"
)

const (
	topicsCount  = 100
	subsPerTopic = 10

	eventsToSend = 10_000
)

func runBenchmark(topics []string) int64 {
	bus := eventbus.New()

	var total int64

	for _, topic := range topics {
		for i := 0; i < subsPerTopic; i++ {
			eventbus.Subscribe(bus, topic, func(event any) bool {
				atomic.AddInt64(&total, 1)
				return false
			})
		}
	}

	for i := 0; i < eventsToSend; i++ {
		index := rand.Intn(len(topics))
		bus.Publish(topics[index], 0)
	}

	bus.Stop()

	return total
}

func BenchmarkEventBus(b *testing.B) {
	topics := make([]string, topicsCount)

	for i := 0; i < topicsCount; i++ {
		topics[i] = uuid.NewString()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if runBenchmark(topics) != 100_000 {
			b.Fail()
		}
	}
}

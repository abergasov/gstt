package network_test

import (
	"sync"
	"testing"

	"github.com/abergasov/gstt/network"
	"github.com/abergasov/gstt/network/singlemap"
	"github.com/abergasov/gstt/network/syncmap"
	"github.com/stretchr/testify/assert"
)

const (
	benchLength = 100_000
)

func runAddBenc(b *testing.B, mt network.MessageTracker) {
	var wg sync.WaitGroup
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func(j int) {
			assert.NoError(b, mt.Add(generateMessage(j)))
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func runDeleteBenc(b *testing.B, mt network.MessageTracker) {
	for i := 0; i < benchLength; i++ {
		assert.NoError(b, mt.Add(generateMessage(i)))
	}
	var wg sync.WaitGroup
	wg.Add(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(j int) {
			// deletes existing and non-existing messages
			_ = mt.Delete(generateMessage(j).ID)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func runMessageBenc(b *testing.B, mt network.MessageTracker) {
	for i := 0; i < benchLength; i++ {
		assert.NoError(b, mt.Add(generateMessage(i)))
	}

	var wg sync.WaitGroup
	wg.Add(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(j int) {
			// gets existing and non-existing messages
			_, _ = mt.Message(generateMessage(j).ID)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func runMessagesBenc(b *testing.B, mt network.MessageTracker) {
	for i := 0; i < benchLength; i++ {
		assert.NoError(b, mt.Add(generateMessage(i)))
	}

	var wg sync.WaitGroup
	wg.Add(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			// gets existing and non-existing messages
			_ = mt.Messages()
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkSinlemapTracker_Add(b *testing.B) {
	mt := singlemap.NewMessageTracker(benchLength)
	runAddBenc(b, mt)
}

func BenchmarkSinlemapTracker_Delete(b *testing.B) {
	mt := singlemap.NewMessageTracker(benchLength)
	runDeleteBenc(b, mt)
}

func BenchmarkSinlemapTracker_Message(b *testing.B) {
	mt := singlemap.NewMessageTracker(benchLength)
	runMessageBenc(b, mt)
}

func BenchmarkSinlemapTracker_Messages(b *testing.B) {
	mt := singlemap.NewMessageTracker(benchLength)
	runMessagesBenc(b, mt)
}

// another implementation of MessageTracker

func BenchmarkMultiplymapTracker_Add(b *testing.B) {
	mt := syncmap.NewMessageTracker(benchLength)
	runAddBenc(b, mt)
}

func BenchmarkMultiplymapTracker_Delete(b *testing.B) {
	mt := syncmap.NewMessageTracker(benchLength)
	runDeleteBenc(b, mt)
}

func BenchmarkMultiplymapTracker_Message(b *testing.B) {
	mt := syncmap.NewMessageTracker(benchLength)
	runMessageBenc(b, mt)
}

func BenchmarkMultiplymapTracker_Messages(b *testing.B) {
	mt := syncmap.NewMessageTracker(benchLength)
	runMessagesBenc(b, mt)
}

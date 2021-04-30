package uid

import (
	"sync"
	"testing"
	"time"
)

// sec := 216485561 / 1e9
//t.Errorf(time.Unix(int64(sec), 1617206400000000000).Format(time.RFC3339))
func TestDefaultUid_NextID(t *testing.T) {
	ig := NewDefaultID(defaultWorker(), defaultStartTimeNano)
	var wg sync.WaitGroup
	wg.Add(100)                        //using 100 goroutine to generate 10000 ids
	results := make(chan int64, 10000) //store result
	for i := 0; i < 100; i++ {
		go func() {
			for i := 0; i < 100; i++ {
				id := ig.NextID()
				t.Logf("id: %d \t %s", id, ParseID(id))
				results <- id
			}
			wg.Done()
		}()
	}
	wg.Wait()

	m := make(map[int64]bool)
	for i := 0; i < 10000; i++ {
		select {
		case id := <-results:
			if _, ok := m[id]; ok {
				t.Errorf("Found duplicated id: %x", id)
			} else {
				m[id] = true
			}
		case <-time.After(2 * time.Second):
			t.Errorf("Expect 10000 ids in results, but got %d", i)
			return
		}
	}
}

func BenchmarkGenID(b *testing.B) {
	ig := NewDefaultID(defaultWorker(), defaultStartTimeNano)
	for i := 0; i < b.N; i++ {
		ig.NextID()
	}
}

func BenchmarkGenIDP(b *testing.B) {
	ig := NewDefaultID(defaultWorker(), defaultStartTimeNano)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ig.NextID()
		}
	})
}

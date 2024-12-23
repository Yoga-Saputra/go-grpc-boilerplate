package kemu

import "testing"

func BenchmarkMutex(b *testing.B) {
	m := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Lock(i).Unlock()
	}
}

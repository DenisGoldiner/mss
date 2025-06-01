package benchmarks

import (
	"fmt"
	"testing"
)

func BenchmarkSample(b *testing.B) {
	// TODO: setup the system
	// TODO: setup inputs

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if x := fmt.Sprintf("%d", 42); x != "42" {
			b.Fatalf("Unexpected string: %s", x)
		}
	}
}

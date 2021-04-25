package wallet

import "testing"

func BenchmarkRegular(b *testing.B) {
	want := int64(2000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := Regular()
		b.StopTimer()
		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer()
	}
}
func BenchmarkConcurrently(b *testing.B) {
	want := int64(2000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := Concurrently()
		b.StopTimer()
		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer()
	}
}

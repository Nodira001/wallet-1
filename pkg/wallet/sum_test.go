package wallet

import (
	"github.com/iqbol007/wallet/pkg/types"
	"testing"
)

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
func BenchmarkService_SumPayments_Single(b *testing.B) {
	s := &Service{}

	s.payments = append(s.payments, &types.Payment{
		ID: "1", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	s.payments = append(s.payments, &types.Payment{
		ID: "12", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	want := types.Money(20_000_00)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := s.SumPayments(1)
		b.StartTimer()
		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer()
	}
}

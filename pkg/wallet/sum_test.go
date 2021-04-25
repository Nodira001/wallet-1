package wallet

import (
	"reflect"
	"testing"

	"github.com/iqbol007/wallet/pkg/types"
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

// func BenchmarkService_SumPayments_Concurrently(b *testing.B) {
// 	s := &Service{}

// 	s.payments = append(s.payments, &types.Payment{
// 		ID: "1", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
// 	})
// 	s.payments = append(s.payments, &types.Payment{
// 		ID: "12", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
// 	})
// 	want := types.Money(20_000_00)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		result := s.SumPayments(2)
// 		b.StartTimer()
// 		if result != want {
// 			b.Fatalf("invalid result, got %v, want %v", result, want)
// 		}
// 		b.StartTimer()
// 	}
// }

func BenchmarkServiceFilter(b *testing.B) {
	s := &Service{}

	s.payments = append(s.payments, &types.Payment{
		ID: "1", AccountID: 11, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	s.payments = append(s.payments, &types.Payment{
		ID: "12", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	s.payments = append(s.payments, &types.Payment{
		ID: "1", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	s.payments = append(s.payments, &types.Payment{
		ID: "12", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	s.payments = append(s.payments, &types.Payment{
		ID: "1", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	s.payments = append(s.payments, &types.Payment{
		ID: "12", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	s.payments = append(s.payments, &types.Payment{
		ID: "1", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	s.payments = append(s.payments, &types.Payment{
		ID: "12", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	s.payments = append(s.payments, &types.Payment{
		ID: "1", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	s.payments = append(s.payments, &types.Payment{
		ID: "12", AccountID: 12, Amount: 10_000_00, Category: "auto", Status: "auto",
	})
	s.payments = append(s.payments, &types.Payment{
		ID: "1", AccountID: 11, Amount: 10_000_00, Category: "auto", Status: "auto",
	})

	want := []types.Payment{{
		ID: "1", AccountID: 11, Amount: 10_000_00, Category: "auto", Status: "auto",
	}, {
		ID: "1", AccountID: 11, Amount: 10_000_00, Category: "auto", Status: "auto",
	}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := s.FilterPayments(11, 3)
		if err != nil {
			b.Fatalf("invalid result, got %v ", result)
		}
		b.StartTimer()
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer()
	}
}

package main

import (
	"fmt"
	"github.com/iqbol007/wallet/pkg/wallet"
)

func main() {

	s := &wallet.Service{}
	s.RegisterAccount("+1")
	s.Deposit(1, 10_000_000_000)
	s.Pay(1, 20, "auto")
	s.Pay(2, 20, "auto")
	s.Pay(3, 20, "auto")
	s.Pay(1, 20, "auto")
	s.Pay(2, 20, "auto")
	s.Pay(3, 20, "auto")
	s.Pay(1, 20, "auto")
	s.Pay(8, 20, "auto")
	s.Pay(6, 20, "auto")
	s.Pay(1, 20, "auto")
	s.Pay(5, 20, "auto")
	res, err := s.FilterPayments(0, 3)
	if err != nil {
		fmt.Printf("1=%v", err)
	}
	fmt.Println(res)

}

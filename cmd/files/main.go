package main

import (
	"github.com/iqbol007/wallet/pkg/wallet"
	"log"
)

func main() {
	s := &wallet.Service{}
	err := s.Import("data")
	log.Print(err)
	////s := &Service{}
	//
	//acc, err := s.RegisterAccount("+992004403883")
	//if err != nil {
	//	fmt.Print(err)
	//	return
	//}
	//
	//err = s.Deposit(acc.ID, 10_000_00)
	//if err != nil {
	//	fmt.Print(err)
	//	return
	//}
	//
	//payment, err := s.Pay(acc.ID, 10_000, "auto")
	//if err != nil {
	//	fmt.Print(err)
	//	return
	//}
	//
	//_, err = s.FavoritePayment(payment.ID, "Auto")
	//if err != nil {
	//	fmt.Print(err)
	//	return
	//}
	//
	//err = s.Export("data")
	//if err != nil {
	//	fmt.Print(err)
	//	return
	//}
}

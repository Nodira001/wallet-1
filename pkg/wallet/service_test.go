package wallet

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestService_RegisterAccount_success(t *testing.T) {
	s := &Service{}
	_, err := s.RegisterAccount("+992004403883")
	if err != nil {
		t.Error("Error")
		return
	}

}
func TestService_RegisterAccount_fail(t *testing.T) {
	s := &Service{}
	s.RegisterAccount("+992004403883")
	_, err := s.RegisterAccount("+992004403883")
	if err == nil {
		t.Error("Error")
		return
	}
}
func TestService_FindAccountByID_success(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	fondedAccount, err := s.FindAccountByID(acc.ID)
	if err != nil {
		t.Error("Error")
		return
	}
	if !reflect.DeepEqual(acc, fondedAccount) {
		t.Error("Error")
		return
	}
}
func TestService_FindAccountByID_fail(t *testing.T) {
	s := &Service{}
	_, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.FindAccountByID(3)
	if err == nil {
		t.Error("Error")
		return
	}

}
func TestService_Deposit_success(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
}
func TestService_Deposit_fail(t *testing.T) {
	s := &Service{}
	err := s.Deposit(12, -10_000_00)
	if err == nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(12, 10_000_00)
	if err == nil {
		t.Error("Error")
		return
	}
}
func TestService_Pay_success(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
}
func TestService_Pay_fail(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.Pay(acc.ID, -1_000_00, "cars")
	if err != ErrAmountMustBePositive {
		t.Error("Error")
		return
	}
	_, err = s.Pay(12, 1_000_00, "cars")
	if err != ErrAccountNotFound {
		t.Error("Error")
		return
	}
	acc.Balance -= 10_000_00
	_, err = s.Pay(acc.ID, 1_000_00, "cars")
	if err != ErrNotEnoughBalance {
		t.Error("Error", err)
		return
	}
}
func TestService_FindPaymentByID_success(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	payment, err := s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Error("Error", err)
		return
	}

}
func TestService_FindPaymentByID_fail(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.FindPaymentByID("12")
	if err != ErrPaymentNotFound {
		t.Error("Error")
		return
	}
}
func TestService_Reject_success(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}

	payment, err := s.Pay(acc.ID, 1_000_00, "cars")

	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Reject(payment.ID)
	if err != nil {
		t.Error("Error")
		return
	}

}
func TestService_Reject_fail(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	payment, err := s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Reject("id")
	if err != ErrPaymentNotFound {
		t.Error("Error")
		return
	}
	payment.AccountID = 22
	err = s.Reject(payment.ID)
	if err != ErrAccountNotFound {
		t.Error("Error")
		return
	}
}
func TestService_Repeat_success(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	payment, err := s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.Repeat(payment.ID)
	if err != nil {
		t.Error("Error")
		return
	}
}
func TestService_Repeat_fail(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.Repeat("id")
	if err != ErrPaymentNotFound {
		t.Error("Error")
		return
	}
}

func TestService_FavoritePayment_success(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	payment, err := s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.FavoritePayment(payment.ID, "Jenya")
	if err != nil {
		t.Error("Error")
		return
	}

}
func TestService_FavoritePayment_fail(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.FavoritePayment("id", "Jenya")
	if err == nil {
		t.Error("Error")
		return
	}
}
func TestSerivice_FindFavoritePaymentByID_success(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	payment, err := s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
	favorite, err := s.FavoritePayment(payment.ID, "Jenya")
	if err != nil {
		t.Error("Error")
		return
	}
	foundedFavorite, err := s.FindFavoriteByID(favorite.ID)
	if err != nil {
		t.Error("Error")
		return
	}
	if !reflect.DeepEqual(favorite, foundedFavorite) {
		t.Error("Error")
		return
	}
}
func TestSerivice_FindFavoritePaymentByID_fail(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	payment, err := s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.FavoritePayment(payment.ID, "Jenya")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.FindFavoriteByID("id")
	if err != ErrFavoriteNotFound {
		t.Error("Error")
		return
	}

}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	payment, err := s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
	favorite, err := s.FavoritePayment(payment.ID, "Jenya")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Error("Error")
		return
	}
}
func TestService_PayFromFavorite_fail(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992004403881")
	if err != nil {
		t.Error("Error")
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Error("Error")
		return
	}
	payment, err := s.Pay(acc.ID, 1_000_00, "cars")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.FavoritePayment(payment.ID, "Jenya")
	if err != nil {
		t.Error("Error")
		return
	}
	_, err = s.PayFromFavorite("id")
	if err != ErrFavoriteNotFound {
		t.Error("Error")
		return
	}
}
func TestService_FullExport(t *testing.T) {
	s := &Service{}

	acc, err := s.RegisterAccount("+992004403883")
	if err != nil {
		fmt.Print(err)
		return
	}

	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		fmt.Print(err)
		return
	}

	payment, err := s.Pay(acc.ID, 10_000, "auto")
	if err != nil {
		fmt.Print(err)
		return
	}

	_, err = s.FavoritePayment(payment.ID, "Auto")
	if err != nil {
		fmt.Print(err)
		return
	}

	err = s.Export("data/")
	if err != nil {
		fmt.Print(err)
		return
	}
}
func TestService_FullImport(t *testing.T) {
	s := &Service{}

	err := s.Import("data/")
	if err != nil {
		t.Error(err)
		return
	}
	log.Print(s.accounts)
}

package wallet

import (
	"github.com/google/uuid"
	"reflect"
	"testing"

	"github.com/iqbol007/wallet/pkg/types"
)

func TestService_FindAccountByID(t *testing.T) {
	svc := &Service{}

	svc.RegisterAccount("+992 004 40 38 83")
	svc.RegisterAccount("+992 004 40 22 11")
	svc.RegisterAccount("+992 004 40 11 22")

	account, err := svc.FindAccountByID(2)

	if err != nil {
		t.Error(err)
		return
	}
	expected := &types.Account{ID: 2, Phone: "+992 004 40 22 11", Balance: 0}
	if !reflect.DeepEqual(expected, account) {
		t.Error(ErrAccountNotFound)
		return
	}

}

func TestService_Reject(t *testing.T) {
	svc := &Service{}

	_, err := svc.RegisterAccount("+992 004 40 38 83")

	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Deposit(1, 1_000)

	if err != nil {
		t.Error(err)
		return
	}

	payment, err := svc.Pay(1, 500, "mobile")

	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Reject(payment.ID)

	if err != nil {
		t.Error(err)
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	s := &Service{}

	account, err := s.addAccountWithBalance("+992 004 40 38 83", 10_000_00)

	if err != nil {
		t.Error(err)
		return
	}

	payment, err := s.Pay(account.ID, 1_000_00, "auto")

	if err != nil {
		t.Error(err)
		return
	}

	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(got, payment) {
		t.Errorf("FindPaymentByID(): wrong returned = %v", err)
		return
	}

}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := &Service{}

	account, err := s.addAccountWithBalance("+992 004 40 38 83", 10_000_00)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.Pay(account.ID, 1_000_00, "auto")

	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {

		t.Error("FindPaymentByID(): must returned error, returned nil")
		return
	}
	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrorNotFound, returned = %v", err)
		return
	}

}

func TestService_Repeat_success(t *testing.T) {
	var svc Service

	account, err := svc.RegisterAccount("+992000000001")

	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	payment, err := svc.Pay(account.ID, 10_00, "Cafe")

	if err != nil {
		t.Errorf("method Pay returned not nil error, account => %v", account)
	}

	pay, err := svc.FindPaymentByID(payment.ID)

	if err != nil {
		t.Errorf("method FindPaymentByID returned not nil error, payment => %v", payment)
	}

	paymentNew, err := svc.Repeat(pay.ID)

	if err != nil {
		t.Errorf("method Repat returned not nil error, paymentNew => %v", paymentNew)
	}

}
func TestService_Repeat_fail(t *testing.T) {
	var s Service

	_, payments, err := s.addAccount(testAccount{
		phone:   "+992 00440 3883",
		balance: 10_000_00,
		payments: []struct {
			amount   types.Money
			category types.PaymentCategory
		}{
			{amount: 1_000_00, category: "auto"},
		},
	})
	if err != nil {
		t.Error("Repeat() err 126")
		return
	}
	payment := payments[0]
	payment.Amount += 10000000
	_, err = s.Repeat(payment.ID)
	if err != ErrNotEnoughBalance {
		t.Error("Repeat() err 132", err)
		return
	}
}

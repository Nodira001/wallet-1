package wallet

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iqbol007/wallet/pkg/types"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be more than 0")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("balance is not enough")
var ErrPaymentNotFound = errors.New("payment is not found")
var ErrFavoriteNotFound = errors.New("ErrFavoriteNotFound")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil

}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {

	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount

	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	return account, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {

	var payment *types.Payment

	for _, pm := range s.payments {
		if pm.ID == paymentID {
			payment = pm
			break
		}
	}

	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

func (s *Service) FindFavoritePaymentByID(paymentID string) (*types.Favorite, error) {

	var payment *types.Favorite

	for _, pm := range s.favorites {
		if pm.ID == paymentID {
			payment = pm
			break
		}
	}

	if payment == nil {
		return nil, ErrFavoriteNotFound
	}
	return payment, nil
}
func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return ErrPaymentNotFound
	}

	account, err := s.FindAccountByID(payment.AccountID)

	if err != nil {
		return ErrAccountNotFound
	}

	account.Balance += payment.Amount
	payment.Status = types.PaymentStatusFail

	return nil
}

func (s *Service) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := s.RegisterAccount(phone)

	if err != nil {
		return nil, err
	}

	err = s.Deposit(account.ID, balance)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {

	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	paymentNew, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}
	return paymentNew, nil
}

func (s *Service) addAccount(data struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("cant register addAcount()")
	}
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("cant deposit addAcount()")
	}
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("cant make payment addAcount()")
		}
	}
	return account, payments, nil
}
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	newFaforite := &types.Favorite{
		ID:        uuid.NewString(),
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category}
	s.favorites = append(s.favorites, newFaforite)
	return newFaforite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favoritePayment, err := s.FindFavoritePaymentByID(favoriteID)
	if err != nil {
		return nil, err
	}
	acc, err := s.FindAccountByID(favoritePayment.AccountID)
	if err != nil {
		return nil, err
	}
	if acc.Balance < favoritePayment.Amount {
		return nil, ErrNotEnoughBalance
	}

	newPayment := types.Payment{
		ID:        uuid.NewString(),
		AccountID: favoritePayment.AccountID,
		Amount:    favoritePayment.Amount,
		Category:  favoritePayment.Category,
		Status:    types.PaymentStatusInProgress,
	}
	acc.Balance -= favoritePayment.Amount
	s.payments = append(s.payments, &newPayment)
	return &newPayment, nil
}

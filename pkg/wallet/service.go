package wallet

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/iqbol007/wallet/pkg/types"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

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

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return nil, ErrAccountNotFound
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return ErrAccountNotFound
	}

	// зачисление средств пока не рассматриваем как платёж
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

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	return s.Pay(payment.AccountID, payment.Amount, payment.Category)
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Name:      name,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}

	return nil, ErrFavoriteNotFound
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
}
func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	for _, account := range s.accounts {
		row := strconv.FormatInt(account.ID, 10) + ";" + string(account.Phone) + ";" + strconv.FormatInt(int64(account.Balance), 10)
		_, err = file.Write([]byte(row))
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}
func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	buf := make([]byte, 1)
	content := make([]byte, 0)
	for {
		read, err := file.Read(buf)

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}
	strings.Split(string(content), ";")
	return nil
}
func (s *Service) Export(dir string) error {
	if len(s.accounts) > 0 {

		accountsFile, err := os.Create(dir + "/accounts.dump")

		if err != nil {
			log.Println(err)
			return err
		}

		defer func() {
			if accErr := accountsFile.Close(); accErr != nil {
				log.Print(accErr)
				return
			}
		}()

		for _, account := range s.accounts {

			accountsRow := strconv.FormatInt(account.ID, 10) + ";" + string(account.Phone) + ";" + strconv.FormatInt(int64(account.Balance), 10)
			_, err = accountsFile.Write([]byte(accountsRow))
			if err != nil {
				log.Print(err)
				return err
			}
		}
	}
	if len(s.payments) > 0 {
		paymentsFile, err := os.Create(dir + "/payments.dump")
		if err != nil {
			log.Println(err)
			return err
		}

		defer func() {
			if payErr := paymentsFile.Close(); payErr != nil {
				log.Print(payErr)
				return
			}
		}()

		for _, payment := range s.payments {

			paymentsRow := payment.ID + ";" + strconv.FormatInt(payment.AccountID, 10) + ";" + strconv.FormatInt(int64(payment.Amount), 10) + ";" + string(payment.Category) + ";" + string(payment.Status) + "\n"
			_, err = paymentsFile.Write([]byte(paymentsRow))
			if err != nil {
				log.Print(err)
				return err
			}
		}
	}

	if len(s.favorites) > 0 {
		favoritesFile, err := os.Create(dir + "/favorites.dump")

		if err != nil {
			log.Println(err)
			return err
		}

		defer func() {
			if favErr := favoritesFile.Close(); favErr != nil {
				log.Print(favErr)
				return
			}
		}()

		for _, favorite := range s.favorites {

			favoriteRow := favorite.ID + ";" + strconv.FormatInt(int64(favorite.AccountID), 10) + ";" + strconv.FormatInt(int64(favorite.Amount), 10) + ";" + favorite.Name + ";" + string(favorite.Category) + "\n"
			_, err = favoritesFile.Write([]byte(favoriteRow))
			if err != nil {
				log.Print(err)
				return err
			}
		}
	}

	return nil
}

func (s *Service) Import(dir string) (importError error) {
	accountsFile, err := os.Open(dir + "/accounts.dump")
	fmt.Println(accountsFile)
	if err == nil {
		defer accountsFile.Close()
		accountsReader := bufio.NewReader(accountsFile)
		for {
			line, err := accountsReader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			account := strings.Split(line, ";")
			accountID, err := strconv.ParseInt(account[0], 10, 64)
			if err != nil {
				return err
			}
			accountPhone := account[1]
			accountBalance, err := strconv.ParseInt(account[2], 10, 64)
			if err != nil {
				return err
			}
			accountBackUp := &types.Account{
				ID:      accountID,
				Phone:   types.Phone(accountPhone),
				Balance: types.Money(accountBalance),
			}
			existent, err := s.FindAccountByID(accountID)
			log.Print(existent)
			log.Println(accountBackUp)
			if err == ErrAccountNotFound {
				s.accounts = append(s.accounts, accountBackUp)
			}
			if !reflect.DeepEqual(existent, accountBackUp) {
				s.accounts = append(s.accounts, accountBackUp)
			}
		}
	} else {
		return err
	}
	paymentsFile, err := os.Open(dir + "/payments.dump")
	if err == nil {
		defer paymentsFile.Close()
		paymentsReader := bufio.NewReader(paymentsFile)
		for {
			line, err := paymentsReader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			payment := strings.Split(line, ";")
			paymentID := payment[0]
			paymentAccountID, err := strconv.ParseInt(payment[1], 10, 64)
			if err != nil {

				return err
			}
			paymentAmount, err := strconv.ParseInt(payment[2], 10, 64)
			if err != nil {
				importError = err
				return importError
			}
			paymentCategory := payment[3]
			paymentStatus := payment[4]
			paymentBackUp := &types.Payment{
				ID:        paymentID,
				AccountID: paymentAccountID,
				Amount:    types.Money(paymentAmount),
				Category:  types.PaymentCategory(paymentCategory),
				Status:    types.PaymentStatus(paymentStatus),
			}
			existent, err := s.FindPaymentByID(paymentID)
			if err == ErrPaymentNotFound {
				s.payments = append(s.payments, paymentBackUp)
			}
			if !reflect.DeepEqual(existent, paymentBackUp) {
				s.payments = append(s.payments, paymentBackUp)
			}
		}
	} else {
		return err
	}
	favoritesFile, err := os.Open(dir + "/favorites.dump")
	if err == nil {
		defer favoritesFile.Close()
		favoritesReader := bufio.NewReader(favoritesFile)
		favorite := []string{}
		for {
			line, err := favoritesReader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			favorite = strings.Split(line, ";")
			favoriteID := favorite[0]
			favoriteAccountID, err := strconv.ParseInt(favorite[1], 10, 64)
			if err != nil {
				return err
			}
			favoriteAmount, err := strconv.ParseInt(favorite[2], 10, 64)
			if err != nil {
				return err
			}
			favoriteName := favorite[3]
			favoriteCategory := favorite[4]
			favoriteBackUp := &types.Favorite{
				ID:        favoriteID,
				AccountID: favoriteAccountID,
				Amount:    types.Money(favoriteAmount),
				Name:      favoriteName,
				Category:  types.PaymentCategory(favoriteCategory),
			}
			existent, err := s.FindFavoriteByID(favoriteID)
			if err == ErrAccountNotFound {
				s.favorites = append(s.favorites, favoriteBackUp)
			}
			if !reflect.DeepEqual(existent, favoriteBackUp) {
				s.favorites = append(s.favorites, favoriteBackUp)
			}
		}
	} else {
		return err
	}
	fmt.Println(s.accounts)
	return nil
}

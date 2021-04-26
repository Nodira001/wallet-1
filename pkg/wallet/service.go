package wallet

import (
	"bufio"
	"errors"
	"github.com/google/uuid"
	"github.com/iqbol007/wallet/pkg/types"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
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
// func (s *Service) ExportToFile(path string) error {
// 	if len(s.accounts) > 0 {
// 		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		var str string
// 		for _, v := range s.accounts {
// 			str += fmt.Sprint(v.ID) + ";" + string(v.Phone) + ";" + fmt.Sprint(v.Balance) + "|"
// 		}
// 		_, err = file.WriteString(str)

// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	}
// 	return nil
// }
// func (s *Service) ImportFromFile(path string) error {
// 	file, err := os.Open(path)
// 	if err != nil {

// 		return err
// 	}
// 	defer file.Close()
// 	buf := make([]byte, 1)
// 	content := make([]byte, 0)
// 	for {
// 		read, err := file.Read(buf)

// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		content = append(content, buf[:read]...)
// 	}
// 	data := strings.Split(string(content), "|")
// 	var row []string
// 	for _, v := range data {
// 		row = append(row, strings.ReplaceAll(v, ";", " "))
// 	}
// 	for _, acc := range row {
// 		if len(acc) == 0 {
// 			continue
// 		}

// 		accountSplit := strings.Split(acc, " ")
// 		id, err := strconv.ParseInt(accountSplit[0], 10, 64)
// 		if err != nil {
// 			return err
// 		}
// 		balance, err := strconv.ParseInt(accountSplit[2], 10, 64)

// 		if err != nil {
// 			return err
// 		}
// 		account := &types.Account{ID: id, Balance: types.Money(balance), Phone: types.Phone(accountSplit[1])}
// 		s.accounts = append(s.accounts, account)

// 	}
// 	return nil
// }

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

			accountsRow := strconv.FormatInt(account.ID, 10) + ";" + string(account.Phone) + ";" + string(strconv.FormatInt(int64(account.Balance), 10)) + "\n"
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

	_, err := os.Stat(dir + "/accounts.dump")

	if err == nil {
		accountsFile, err := os.Open(dir + "/accounts.dump")
		if err != nil {
			return err
		}
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
			accountBalance, err := strconv.ParseInt(strings.ReplaceAll(account[2], "\n", ""), 10, 64)
			if err != nil {
				return err
			}
			accountBackUp := &types.Account{
				ID:      accountID,
				Phone:   types.Phone(accountPhone),
				Balance: types.Money(accountBalance),
			}
			_, err = s.FindAccountByID(accountID)
			if err == ErrAccountNotFound {
				s.accounts = append(s.accounts, accountBackUp)
				s.nextAccountID = int64(len(s.accounts))
			}

		}
	}
	_, err = os.Stat(dir + "/payments.dump")
	if err == nil {
		paymentsFile, err := os.Open(dir + "/payments.dump")
		if err != nil {
			return err
		}

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
			paymentStatus := strings.ReplaceAll(payment[4], "\n", "")
			paymentBackUp := &types.Payment{
				ID:        paymentID,
				AccountID: paymentAccountID,
				Amount:    types.Money(paymentAmount),
				Category:  types.PaymentCategory(paymentCategory),
				Status:    types.PaymentStatus(paymentStatus),
			}
			_, err = s.FindPaymentByID(paymentID)
			if err == ErrPaymentNotFound {
				s.payments = append(s.payments, paymentBackUp)
			}

		}
	}

	_, err = os.Stat(dir + "/favorites.dump")
	if err == nil {
		favoritesFile, err := os.Open(dir + "/favorites.dump")
		if err != nil {
			return err
		}
		defer favoritesFile.Close()
		favoritesReader := bufio.NewReader(favoritesFile)

		for {
			line, err := favoritesReader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			favorite := strings.Split(line, ";")
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
			favoriteCategory := strings.ReplaceAll(favorite[4], "\n", "")
			favoriteBackUp := &types.Favorite{
				ID:        favoriteID,
				AccountID: favoriteAccountID,
				Amount:    types.Money(favoriteAmount),
				Name:      favoriteName,
				Category:  types.PaymentCategory(favoriteCategory),
			}
			_, err = s.FindFavoriteByID(favoriteID)
			if err == ErrFavoriteNotFound {
				s.favorites = append(s.favorites, favoriteBackUp)
			}

		}
	}

	return nil
}

// func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
// 	var payments []types.Payment
// 	acc, err := s.FindAccountByID(accountID)
// 	if err != nil {
// 		return nil, ErrAccountNotFound
// 	}
// 	for _, payment := range s.payments {
// 		if acc.ID == payment.AccountID {
// 			payments = append(payments, *payment)
// 		}
// 	}
// 	return payments, nil
// }

// func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
// 	if len(payments) > 0 {
// 		if len(payments) <= records {
// 			file, _ := os.OpenFile(dir+"/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
// 			defer file.Close()

// 			var str string
// 			for _, v := range payments {
// 				str += fmt.Sprint(v.ID) + ";" + fmt.Sprint(v.AccountID) + ";" + fmt.Sprint(v.Amount) + ";" + fmt.Sprint(v.Category) + ";" + fmt.Sprint(v.Status) + "\n"
// 			}
// 			file.WriteString(str)
// 		} else {

// 			var row string
// 			k := 0
// 			count := 1
// 			var file *os.File
// 			for _, v := range payments {
// 				if k == 0 {
// 					file, _ = os.OpenFile(dir+"/payments"+fmt.Sprint(count)+".dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
// 				}
// 				k++
// 				row = fmt.Sprint(v.ID) + ";" + fmt.Sprint(v.AccountID) + ";" + fmt.Sprint(v.Amount) + ";" + fmt.Sprint(v.Category) + ";" + fmt.Sprint(v.Status) + "\n"
// 				_, err := file.WriteString(row)
// 				if err != nil {
// 					return err
// 				}
// 				if k == records {
// 					row = ""
// 					count++
// 					k = 0
// 					file.Close()
// 				}
// 			}

// 		}
// 	}
// 	return nil
// }

func (s *Service) SumPayments(goroutines int) types.Money {

	mu := sync.Mutex{}
	sum := types.Money(0)

	if goroutines == 0 || goroutines == 1 {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			val := types.Money(0)
			for _, payment := range s.payments {
				val += payment.Amount
			}
			mu.Lock()
			defer mu.Unlock()
			sum += val
		}()
		wg.Wait()
		return sum
	}
	wg := sync.WaitGroup{}
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {

		go func() {
			defer wg.Done()
			val := types.Money(0)
			for _, payment := range s.payments {
				val += payment.Amount
			}
			mu.Lock()
			defer mu.Unlock()
			sum += val
		}()
	}

	wg.Wait()
	return sum / 10
}

func (s *Service) FilterPaymentsForGoroutines(goroutinesCount int, accountID int64) ([][]types.Payment, error) {
	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	payments := []types.Payment{}

	for _, p := range s.payments {

		if p.AccountID == accountID {

			payments = append(payments, *p)

		}
	}

	grouped := [][]types.Payment{}

	for i := 0; i < len(payments); i++ {

		if i+goroutinesCount > len(payments)-1 {

			grouped = append(grouped, payments[i:])

			break
		}

		grouped = append(grouped, payments[i:i+goroutinesCount])

		i += goroutinesCount - 1
	}

	return grouped, nil
}

func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {
	if goroutines == 0 {
		mu := sync.Mutex{}
		payments := []types.Payment{}

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			val := []types.Payment{}
			for _, payment := range s.payments {
				if accountID == payment.AccountID {
					val = append(payments, *payment)
				}
			}
			mu.Lock()
			defer mu.Unlock()
			payments = append(payments, val...)

		}()

		wg.Wait()
		if len(payments) == 0 {

			return nil, ErrAccountNotFound
		}

		return payments, nil
	}

	wg := sync.WaitGroup{}

	mu := sync.Mutex{}
	payments := []types.Payment{}

	filteredPayments, err := s.FilterPaymentsForGoroutines(goroutines, accountID)
	if err != nil {
		return nil, err
	}
	if len(filteredPayments) == 0 {
		return nil, nil
	}
	for _, fp := range filteredPayments {
		wg.Add(1)
		go func(fp []types.Payment) {
			defer wg.Done()
			mu.Lock()
			payments = append(payments, fp...)
			defer mu.Unlock()
		}(fp)
	}

	wg.Wait()
	if len(payments) == 0 {
		return nil, nil
	}

	return payments, nil

}

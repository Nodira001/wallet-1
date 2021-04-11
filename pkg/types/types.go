package types

type Money int64

type PaymentCategory string

type PaymentStatus string

const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "InProgress"
)

type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

type Phone string

type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

type Messenger interface {
	Send(message string) (ok bool)
	Receive() (message string, ok bool)
}

type Favorite struct {
	ID        string
	AccountID int64
	Name      string
	Amount    Money
	Category  PaymentCategory
}

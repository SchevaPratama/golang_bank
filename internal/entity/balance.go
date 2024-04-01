package entity

type Balance struct {
	ID       string
	UserId   string
	Balance  int
	Currency string
}

func (prod *Balance) TableName() string {
	return "balances"
}

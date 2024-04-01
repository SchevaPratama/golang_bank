package entity

type Transaction struct {
	ID                 string
	UserId             string
	BankAccountNumber  string
	BankName           string
	TransferProofImage string
	Balance            int
	Currency           string
	CreatedAt          string
	UpdatedAt          string
	Flow               string
}

func (prod *Transaction) TableName() string {
	return "transactions"
}

package converter

import (
	"fmt"
	"log"
	"time"

	"github.com/paimon_bank/internal/entity"
	"github.com/paimon_bank/internal/model"
)

func TransactionConverter(transaction *entity.Transaction) *model.TransactionResponse {
	t, err := time.Parse("2006-01-02T15:04:05.999999Z", transaction.CreatedAt)
	if err != nil {
		log.Println("Error parsing timestamp:", err)
	}

	if transaction.Flow == "transaction" {
		return &model.TransactionResponse{
			TransactionId:    transaction.ID,
			Balance:          -transaction.Balance,
			Currency:         transaction.Currency,
			TransferProofImg: transaction.TransferProofImage,
			CreatedAt:        fmt.Sprintf("%d", t.UnixMilli()),
			Source: map[string]string{
				"bankAccountNumber": transaction.BankAccountNumber,
				"bankName":          transaction.BankName,
			},
		}
	}

	if transaction.Flow == "topup" {
		return &model.TransactionResponse{
			TransactionId:    transaction.ID,
			Balance:          transaction.Balance,
			Currency:         transaction.Currency,
			TransferProofImg: transaction.TransferProofImage,
			CreatedAt:        fmt.Sprintf("%d", t.UnixMilli()),
			Source: map[string]string{
				"bankAccountNumber": transaction.BankAccountNumber,
				"bankName":          transaction.BankName,
			},
		}
	}

	return nil
}

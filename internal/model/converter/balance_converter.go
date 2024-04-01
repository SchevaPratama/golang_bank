package converter

import (
	"github.com/paimon_bank/internal/entity"
	"github.com/paimon_bank/internal/model"
)

func BalanceConverter(balance *entity.Balance) *model.BalanceResponse {
	return &model.BalanceResponse{
		Balance:  balance.Balance,
		Currency: balance.Currency,
	}
}

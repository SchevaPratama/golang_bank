package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ladydascalie/currency"
	"github.com/paimon_bank/internal/customErr"
	"github.com/paimon_bank/internal/entity"
	"github.com/paimon_bank/internal/model"
	"github.com/paimon_bank/internal/model/converter"
	"github.com/paimon_bank/internal/repository"
	"github.com/sirupsen/logrus"
)

type TransactionService struct {
	Repository        *repository.TransactionRepository
	RepositoryBalance *repository.BalanceRepository
	Validate          *validator.Validate
	Log               *logrus.Logger
}

func NewTransactionService(r *repository.TransactionRepository, b *repository.BalanceRepository, log *logrus.Logger, validate *validator.Validate) *TransactionService {
	return &TransactionService{Repository: r, RepositoryBalance: b, Validate: validate, Log: log}
}

func (s *TransactionService) Create(ctx context.Context, request *model.TransactionRequest, userId string) error {
	err := request.Validate()
	if err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	if !currency.Valid(request.FromCurrency) {
		custErr := customErr.NewBadRequestError("currency format is not ISO 4217")
		return custErr
	}

	dataBalances, _ := s.RepositoryBalance.GetBalance(userId)

	newRequest := &entity.Transaction{
		ID:                uuid.New().String(),
		UserId:            userId,
		BankAccountNumber: request.RecepientBankAccountNumber,
		BankName:          request.RecepientBankName,
		Balance:           request.Balances,
		Currency:          request.FromCurrency,
	}

	var balance entity.Balance
	for _, blnc := range *dataBalances {
		if blnc.Currency == newRequest.Currency {
			balance = blnc
		}
	}

	if balance.Balance-newRequest.Balance < 0 {
		custErr := customErr.NewBadRequestError("balance is not enough")
		return custErr
	}

	errTrans := s.Repository.CreateTransaction(newRequest)
	if errTrans != nil {
		return customErr.NewInternalServerError(errTrans.Error())
	}

	return nil
}

func (s *TransactionService) TransactionHistory(ctx context.Context, filter *model.TransactionFilter, userId string) ([]model.TransactionResponse, error) {
	var listTransaction *[]entity.Transaction

	listTransaction, err := s.Repository.GetTransactions(filter, userId)
	if err != nil {
		return nil, err
	}

	newTransactions := make([]model.TransactionResponse, len(*listTransaction))
	for i, transaction := range *listTransaction {
		newTransactions[i] = *converter.TransactionConverter(&transaction)
	}
	return newTransactions, err
}

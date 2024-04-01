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

type BalanceService struct {
	Repository *repository.BalanceRepository
	Validate   *validator.Validate
	Log        *logrus.Logger
}

func NewBalanceService(r *repository.BalanceRepository, log *logrus.Logger, validate *validator.Validate) *BalanceService {
	return &BalanceService{Repository: r, Validate: validate, Log: log}
}

func (s *BalanceService) Create(ctx context.Context, request *model.BalanceRequest, userId string) error {
	err := request.Validate()
	if err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	if !currency.Valid(request.Currency) {
		custErr := customErr.NewBadRequestError("currency format is not ISO 4217")
		return custErr
	}

	dataBalances, _ := s.Repository.GetBalance(userId)

	newRequest := &entity.Balance{
		ID:       uuid.New().String(),
		UserId:   userId,
		Balance:  request.AddedBalance,
		Currency: request.Currency,
	}

	var balance entity.Balance
	for _, blnc := range *dataBalances {
		if blnc.Currency == newRequest.Currency {
			balance = blnc
		}
	}

	if balance.ID != "" {
		err := s.Repository.UpdateBalance(newRequest, request)
		if err != nil {
			return customErr.NewInternalServerError(err.Error())
		}
	}

	if balance.ID == "" {
		errRepo := s.Repository.Create(newRequest, request)
		if errRepo != nil {
			return customErr.NewInternalServerError(errRepo.Error())
		}
	}

	return nil
}

func (s *BalanceService) ListBalance(ctx context.Context, userId string) ([]model.BalanceResponse, error) {
	var listBalance *[]entity.Balance

	listBalance, err := s.Repository.GetBalance(userId)
	if err != nil {
		return nil, err
	}

	newBalances := make([]model.BalanceResponse, len(*listBalance))
	for i, balance := range *listBalance {
		newBalances[i] = *converter.BalanceConverter(&balance)
	}
	return newBalances, err
}

package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-playground/validator/v10"
	"github.com/ladydascalie/currency"
	// "github.com/golang-jwt/jwt/v4"
)

type TransactionResponse struct {
	TransactionId    string            `json:"transactionId"`
	Balance          int               `json:"balance"`
	Currency         string            `json:"currency"`
	TransferProofImg string            `json:"transferProofImg"`
	CreatedAt        string            `json:"createdAt"`
	Source           map[string]string `json:"source"`
}

type TransactionRequest struct {
	RecepientBankAccountNumber string `json:"recipientBankAccountNumber" validate:"required,min=5,max=30"`
	RecepientBankName          string `json:"recipientBankName" validate:"required,min=5,max=30"`
	Balances                   int    `json:"balances" validate:"required,min=0"`
	FromCurrency               string `json:"fromCurrency" validate:"required"`
}

type TransactionFilter struct {
	Limit  int `json:"limit" validate:"numeric,min=0"`
	Offset int `json:"offset" validate:"numeric,min=0"`
}

func (b *TransactionRequest) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
		return currency.Valid(fl.Field().String())
	})

	err := validation.ValidateStruct(b,
		validation.Field(&b.RecepientBankAccountNumber,
			validation.Required.Error("bank account number is required"),
			validation.Length(5, 30).Error("bank account number must be between 5 and 30 characters"),
		),
		validation.Field(&b.RecepientBankName,
			validation.Required.Error("bank account name is required"),
			validation.Length(5, 30).Error("bank account name must be between 5 and 30 characters"),
		),
		validation.Field(&b.Balances,
			validation.Required.Error("added balance is required"),
			validation.Min(0).Error("added balance must be minimal 0"),
		),
		validation.Field(&b.FromCurrency,
			validation.Required.Error("currency is required"),
		),
	)

	return err
}

package model

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-playground/validator/v10"
	"github.com/ladydascalie/currency"
	// "github.com/golang-jwt/jwt/v4"
)

type BalanceResponse struct {
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
}

type BalanceRequest struct {
	SenderBankAccountNumber string `json:"senderBankAccountNumber" validate:"required,min=5,max=30"`
	SenderBankName          string `json:"senderBankName" validate:"required,min=5,max=30"`
	AddedBalance            int    `json:"addedBalance" validate:"required,min=0"`
	Currency                string `json:"currency" validate:"required"`
	TransferProofImage      string `json:"transferProofImg" validate:"required,url"`
}

type ProductFilter struct {
	Condition      *string   `json:"condition"`
	Keyword        *string   `json:"keyword"`
	SortBy         *string   `json:"sortBy"`
	OrderBy        *string   `json:"orderBy"`
	MaxPrice       *int      `json:"maxPrice"`
	MinPrice       *int      `json:"minPrice"`
	Tags           *[]string `json:"tags"`
	UserOnly       *bool     `json:"userOnly"`
	ShowEmptyStock *bool     `json:"showEmptyStock"`
	Limit          *int      `json:"limit"`
	Offset         *int      `json:"offset"`
}

type StockRequest struct {
	Stock int16 `json:"stock" validate:"min=0"`
}

type BuyRequest struct {
	BankAccountId        string `json:"bankAccountId" validate:"required"`
	ProductId            string `json:"productId" validate:"required"`
	PaymentProofImageUrl string `json:"paymentProofImageUrl" validate:"required,url"`
	Quantity             int16  `json:"quantity" validate:"required,min=1"`
}

func (b *BalanceRequest) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
		return currency.Valid(fl.Field().String())
	})

	err := validation.ValidateStruct(b,
		validation.Field(&b.SenderBankAccountNumber,
			validation.Required.Error("bank account number is required"),
			validation.Length(5, 30).Error("bank account number must be between 5 and 30 characters"),
		),
		validation.Field(&b.SenderBankName,
			validation.Required.Error("bank account name is required"),
			validation.Length(5, 30).Error("bank account name must be between 5 and 30 characters"),
		),
		validation.Field(&b.AddedBalance,
			validation.Required.Error("added balance is required"),
			validation.Min(0).Error("added balance must be minimal 0"),
		),
		validation.Field(&b.Currency,
			validation.Required.Error("currency is required"),
		),
		validation.Field(&b.TransferProofImage,
			validation.Required.Error("transfer proof image is required"),
			validation.By(validateImage),
		),
	)

	return err
}

func validateImage(value any) error {
	image, _ := value.(string)

	pattern := `http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+.(?:jpg|jpeg|png|gif|bmp|webp|svg)$`
	rgx := regexp.MustCompile(pattern)
	if !rgx.MatchString(image) {
		return errors.New("invalid image url format")
	}

	return nil
}

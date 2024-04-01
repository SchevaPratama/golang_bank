package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	// "github.com/golang-jwt/jwt/v4"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRegisterResponse struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

func (r *LoginRequest) LoginRequestValidate() error {
	err := validation.ValidateStruct(r,
		validation.Field(&r.Password,
			validation.Required.Error("password is required"),
		),
		validation.Field(&r.Email,
			validation.Required.Error("email is required"),
			validation.By(validateEmailFormat),
		),
	)

	return err

}

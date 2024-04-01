package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	jtoken "github.com/golang-jwt/jwt/v4"

	"github.com/go-playground/validator/v10"
	"github.com/paimon_bank/internal/customErr"
	"github.com/paimon_bank/internal/entity"
	"github.com/paimon_bank/internal/model"
	"github.com/paimon_bank/internal/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repository *repository.UserRepository
	Validate   *validator.Validate
	Log        *logrus.Logger
}

func NewUserService(r *repository.UserRepository, log *logrus.Logger, validate *validator.Validate) *UserService {
	return &UserService{Repository: r, Validate: validate, Log: log}
}

func (s *UserService) Register(ctx context.Context, request *model.RegisterRequest) (*model.LoginRegisterResponse, error) {
	err := request.RegisterValidate()
	if err != nil {
		return nil, customErr.NewBadRequestError(err.Error())
	}

	bcryptSalt, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcryptSalt)
	if err != nil {
		s.Log.WithError(err).Error("Error hashedPassword")
	}

	payload := &entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	user, _ := s.Repository.GetByEmail(payload.Email)
	if user != nil {
		return nil, customErr.NewConflictError("user already exist")
	}

	err = s.Repository.Create(payload)
	if err != nil {
		return nil, err
	}

	day := time.Hour * 8

	claims := jtoken.MapClaims{
		"ID":    payload.ID,
		"email": payload.Email,
		"exp":   time.Now().Add(day * 1).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	resp := &model.LoginRegisterResponse{
		Email:       payload.Email,
		Name:        payload.Name,
		AccessToken: t,
	}

	return resp, nil
}

func (s *UserService) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginRegisterResponse, error) {
	err := request.LoginRequestValidate()
	if err != nil {
		return nil, customErr.NewBadRequestError(err.Error())
	}

	payload := &entity.User{
		Email:    request.Email,
		Password: request.Password,
	}

	user, _ := s.Repository.GetByEmail(payload.Email)
	if user == nil {
		return nil, customErr.NewNotFoundError("user not found")
	}

	fmt.Println(user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		s.Log.WithError(err).Error("Error Password is Wrong", err.Error())
		return nil, customErr.NewBadRequestError("password is wrong")
	}

	day := time.Hour * 8

	claims := jtoken.MapClaims{
		"ID":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(day * 1).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	resp := &model.LoginRegisterResponse{
		Email:       user.Email,
		Name:        user.Name,
		AccessToken: t,
	}

	return resp, nil
}

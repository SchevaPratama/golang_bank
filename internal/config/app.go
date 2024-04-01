package config

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/paimon_bank/internal/handler"
	"github.com/paimon_bank/internal/repository"
	"github.com/paimon_bank/internal/router"
	"github.com/paimon_bank/internal/service"
	"github.com/sirupsen/logrus"
)

type BootstrapConfig struct {
	DB       *sqlx.DB
	App      *fiber.App
	Log      *logrus.Logger
	Aws      *s3.Client
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	balanceRepository := repository.NewBalanceRepository(config.DB)
	userRepository := repository.NewUserRepository(config.DB)
	transactionRepository := repository.NewTransactionRepository(config.DB)

	// setup services
	imageService := service.NewImageService(config.Aws, config.Log)
	balanceService := service.NewBalanceService(balanceRepository, config.Log, config.Validate)
	userService := service.NewUserService(userRepository, config.Log, config.Validate)
	transactionservice := service.NewTransactionService(transactionRepository, balanceRepository, config.Log, config.Validate)

	// setup handler
	ImageHandler := handler.NewImageHandler(imageService, config.Log)
	BalanceHandler := handler.NewBalanceHandler(balanceService, config.Log, config.Validate)
	UserHandler := handler.NewUserHandler(userService, config.Log)
	TransactionHandler := handler.NewTranssactionHandler(transactionservice, config.Log, config.Validate)

	// recover from panic
	config.App.Use(func(c *fiber.Ctx) error {

		defer func() {
			if r := recover(); r != nil {
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal Server Error",
				})
			}
		}()
		return c.Next()
	})

	// setup route
	routeConfig := router.RouteConfig{
		App:                config.App,
		ImageHandler:       ImageHandler,
		BalanceHandler:     BalanceHandler,
		UserHandler:        UserHandler,
		TransactionHandler: TransactionHandler,
	}

	routeConfig.Setup()
}

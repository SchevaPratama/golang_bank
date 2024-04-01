package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paimon_bank/internal/middleware"
)

func NewFiber() *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      "Paimon Bank",
		ErrorHandler: middleware.ErrorHandler,
		Prefork:      false,
	})

	return app
}

// func NewErrorHandler() fiber.ErrorHandler {
// 	return func(ctx *fiber.Ctx, err error) error {
// 		code := fiber.StatusInternalServerError
// 		if e, ok := err.(*fiber.Error); ok {
// 			code = e.Code
// 		}

// 		return ctx.Status(code).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
// }

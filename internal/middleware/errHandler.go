package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paimon_bank/internal/customErr"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := ctx.Response().StatusCode()
	if customError, ok := err.(customErr.CustomError); ok {
		code = customError.Status()
		return ctx.Status(code).JSON(fiber.Map{
			"messages": customError.Error(),
		})
	} else if code < 400 {
		code = fiber.StatusInternalServerError
	}

	return ctx.Status(code).JSON(fiber.Map{
		"messages": err.Error(),
	})
}

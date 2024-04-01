package handler

import (
	// "golang_socmed/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/paimon_bank/internal/customErr"
	"github.com/paimon_bank/internal/model"
	"github.com/paimon_bank/internal/service"
	"github.com/sirupsen/logrus"
)

type BalanceHandler struct {
	Service  *service.BalanceService
	Log      *logrus.Logger
	validate *validator.Validate
}

func NewBalanceHandler(s *service.BalanceService, log *logrus.Logger, validate *validator.Validate) *BalanceHandler {
	return &BalanceHandler{
		Service:  s,
		Log:      log,
		validate: validate,
	}
}

func (h *BalanceHandler) Create(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	request := new(model.BalanceRequest)
	if err := c.BodyParser(&request); err != nil {
		custErr := customErr.NewBadRequestError(err.Error())
		return c.Status(custErr.StatusCode).JSON(fiber.Map{"message": custErr.Message})
	}

	err := h.Service.Create(c.UserContext(), request, userId)
	if err != nil {
		// return fiber.ErrBadRequest
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    1,
		"message": "add balance success",
		"data":    request,
	})
}

func (h *BalanceHandler) ListBalance(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}
	dataBalances, err := h.Service.ListBalance(c.UserContext(), userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    dataBalances,
	})
}

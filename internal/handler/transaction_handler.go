package handler

import (
	// "golang_socmed/internal/model"

	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/paimon_bank/internal/customErr"
	"github.com/paimon_bank/internal/model"
	"github.com/paimon_bank/internal/service"
	"github.com/sirupsen/logrus"
)

type TranssactionHandler struct {
	Service  *service.TransactionService
	Log      *logrus.Logger
	validate *validator.Validate
}

func NewTranssactionHandler(s *service.TransactionService, log *logrus.Logger, validate *validator.Validate) *TranssactionHandler {
	return &TranssactionHandler{
		Service:  s,
		Log:      log,
		validate: validate,
	}
}

func (h *TranssactionHandler) Create(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}

	request := new(model.TransactionRequest)
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
		"message": "transaction success",
		"data":    request,
	})
}

func (h *TranssactionHandler) TransactionHistory(c *fiber.Ctx) error {
	userId, ok := c.Locals("userLoggedInId").(string)
	if !ok {
		return &fiber.Error{
			Code:    500,
			Message: "Failed",
		}
	}
	params := c.Queries()

	filter := model.TransactionFilter{
		Limit:  5,
		Offset: 0,
	}

	if val, ok := params["limit"]; ok {
		if val != "" {
			limitParsed, _ := strconv.Atoi(val)
			if limitParsed < 0 {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"message": "limit minimal 0",
				})
			}
			filter.Limit = limitParsed
		} else {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "limit can't be empty",
			})
		}
	} else {
		filter.Limit = 5
	}

	if val, ok := params["offset"]; ok {
		if val != "" {
			offsetParsed, _ := strconv.Atoi(val)
			if offsetParsed < 0 {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"message": "offset minimal 0",
				})
			}
			filter.Offset = offsetParsed
		} else {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "offset can't be empty",
			})
		}
	} else {
		filter.Offset = 0
	}

	if err := c.QueryParser(&filter); err != nil {
		custErr := customErr.NewBadRequestError(err.Error())
		return c.Status(custErr.StatusCode).JSON(fiber.Map{"message": custErr.Message})
	}

	dataBalances, err := h.Service.TransactionHistory(c.UserContext(), &filter, userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    dataBalances,
	})
}

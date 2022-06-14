package handlers

import (
	"cashier/app/models"
	"cashier/app/services"
	"cashier/utils"

	"github.com/gofiber/fiber"
)

type CashHandler struct {
	cashService services.CashService
}

func NewCashHandler(cashService services.CashService) *CashHandler {
	return &CashHandler{
		cashService: cashService,
	}
}

func (h *CashHandler) GetStore(c *fiber.Ctx) {
	store, err := h.cashService.GetStore()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}
	utils.SuccessResponse(c, store)
}

func (h *CashHandler) ReplaceStore(c *fiber.Ctx) {
	requestModel := new([]models.CashStoreBody)
	if err := c.BodyParser(requestModel); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	err := h.cashService.ReplaceStore(*requestModel)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *CashHandler) Pay(c *fiber.Ctx) {
	requestModel := new(models.PayRequestBody)
	if err := c.BodyParser(requestModel); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	if len(requestModel.PaymentType) == 0 {
		utils.BadRequestResponse(c, "payment type is required")
		return
	}

	if requestModel.PaymentType != models.Cash &&
		requestModel.PaymentType != models.EWallet &&
		requestModel.PaymentType != models.CreditCard {

		utils.BadRequestResponse(c, "payment type is invalid")
		return
	}

	pay, err := h.cashService.Pay(*requestModel)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, pay)
}

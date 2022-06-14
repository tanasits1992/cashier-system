package handlers

import (
	"cashier/app/services"
	"cashier/utils"

	"github.com/gofiber/fiber"
)

type SaleTransactionHandler struct {
	saleTransactionService services.SaleTransactionService
}

func NewSaleTransactionHandler(saleTransactionService services.SaleTransactionService) *SaleTransactionHandler {
	return &SaleTransactionHandler{
		saleTransactionService: saleTransactionService,
	}
}

func (h *SaleTransactionHandler) GetByBillNo(c *fiber.Ctx) {
	billNo := c.Query("billNo")
	if len(billNo) == 0 {
		utils.BadRequestResponse(c, "bill no is required")
		return
	}

	saleTransaction, err := h.saleTransactionService.GetByBillNo(billNo)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}
	utils.SuccessResponse(c, saleTransaction)
}

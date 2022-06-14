package handlers

import (
	"cashier/app/models"
	"cashier/app/services"
	"cashier/utils"

	"github.com/gofiber/fiber"
)

type VoucherHandler struct {
	voucherService services.VoucherService
}

func NewVoucherHandler(voucherService services.VoucherService) *VoucherHandler {
	return &VoucherHandler{
		voucherService: voucherService,
	}
}

func (h *VoucherHandler) Insert(c *fiber.Ctx) {
	requestModel := new(models.CreateVoucherRequestBody)
	if err := c.BodyParser(requestModel); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	if len(requestModel.Name) < 0 {
		utils.BadRequestResponse(c, "name is require")
		return
	}
	if requestModel.Discount < 0 {
		utils.BadRequestResponse(c, "discount is require")
		return
	}
	if requestModel.Start.IsZero() {
		utils.BadRequestResponse(c, "start is require")
		return
	}
	if requestModel.End.IsZero() {
		utils.BadRequestResponse(c, "end is require")
		return
	}

	barcode, err := h.voucherService.Insert(*requestModel)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, fiber.Map{
		"barcode": barcode,
	})
}

func (h *VoucherHandler) Inactivate(c *fiber.Ctx) {
	code := c.Params("code")
	if len(code) == 0 {
		utils.BadRequestResponse(c, "code is require")
		return
	}

	err := h.voucherService.Inactivate(code)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *VoucherHandler) List(c *fiber.Ctx) {
	vouchers, err := h.voucherService.List()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, vouchers)
}

func (h *VoucherHandler) GetByCode(c *fiber.Ctx) {
	code := c.Params("code")
	if len(code) == 0 {
		utils.BadRequestResponse(c, "code is require")
		return
	}

	voucher, err := h.voucherService.GetByCode(code)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, voucher)
}

func (h *VoucherHandler) Validate(c *fiber.Ctx) {
	code := c.Params("code")
	if len(code) == 0 {
		utils.BadRequestResponse(c, "code is require")
		return
	}

	pass, detail, err := h.voucherService.Validate(code)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, fiber.Map{
		"pass":   pass,
		"detail": detail,
	})
}

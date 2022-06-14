package routes

import (
	"cashier/app/handlers"
	"cashier/app/services"
	"cashier/fake_database"

	"github.com/gofiber/fiber"
)

func voucherSetup() *handlers.VoucherHandler {
	voucherRepository := fake_database.NewVoucherRepository()
	voucherservice := services.NewVoucherService(voucherRepository)

	return handlers.NewVoucherHandler(*voucherservice)
}

func VoucherRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	voucherHandler := voucherSetup()

	route.Post("/voucher", voucherHandler.Insert)
	route.Patch("/voucher/:code", voucherHandler.Inactivate)
	route.Get("/voucher", voucherHandler.List)
	route.Get("/voucher/:code", voucherHandler.GetByCode)
	route.Get("/voucher/validate/:code", voucherHandler.Validate)

}

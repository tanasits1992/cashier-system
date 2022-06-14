package routes

import (
	"cashier/app/handlers"
	"cashier/app/services"
	"cashier/fake_database"

	"github.com/gofiber/fiber"
)

func saleTransactionSetup() *handlers.SaleTransactionHandler {
	saleTransactionRepository := fake_database.NewSaleTransactionRepository()
	saleTransactionservice := services.NewSaleTransactionService(saleTransactionRepository)
	return handlers.NewSaleTransactionHandler(*saleTransactionservice)
}

func SaleTransactionRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	saleTransactionHandler := saleTransactionSetup()

	route.Get("/sale-transaction", saleTransactionHandler.GetByBillNo)
}

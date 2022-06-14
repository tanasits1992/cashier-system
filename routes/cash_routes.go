package routes

import (
	"cashier/app/handlers"
	"cashier/app/services"
	"cashier/connectors"
	"cashier/fake_database"

	"github.com/gofiber/fiber"
)

func cashStoreSetup() *handlers.CashHandler {
	eWalletConnector := connectors.NewEWalletConnector()
	creditCardConnector := connectors.NewCerditCardConnector()
	notificationConnector := connectors.NewNotificationConnector()
	cashStoreRepository := fake_database.NewCashStoreRepository()
	voucherRepository := fake_database.NewVoucherRepository()
	saleTransactionRepository := fake_database.NewSaleTransactionRepository()

	voucherService := services.NewVoucherService(voucherRepository)

	cashStoreService := services.NewCashService(
		eWalletConnector,
		creditCardConnector,
		notificationConnector,
		*voucherService,
		cashStoreRepository,
		voucherRepository,
		saleTransactionRepository,
	)

	return handlers.NewCashHandler(*cashStoreService)
}

func CashierRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	cashStoreHandler := cashStoreSetup()

	route.Get("/cash/store", cashStoreHandler.GetStore)
	route.Put("/cash/store", cashStoreHandler.ReplaceStore)
	route.Post("/cash/pay", cashStoreHandler.Pay)
}

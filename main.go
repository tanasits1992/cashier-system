package main

import (
	"cashier/fake_database"
	"cashier/routes"
	"log"

	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	app := fiber.New()

	storeRepository := fake_database.NewCashStoreRepository()
	storeRepository.Load("cash-store.json")

	voucherRepository := fake_database.NewVoucherRepository()
	voucherRepository.Load("voucher.json")

	routes.SeriesRoutes(app)
	routes.VoucherRoutes(app)
	routes.CashierRoutes(app)
	routes.SaleTransactionRoutes(app)

	app.Listen("3000")
}

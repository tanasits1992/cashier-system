package routes

import (
	"cashier/app/handlers"
	"cashier/app/services"

	"github.com/gofiber/fiber"
)

func seriesSetup() *handlers.SeriesHandler {
	seriesService := services.NewSeriesService()

	return handlers.NewSeriesHandler(*seriesService)
}

func SeriesRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	seriesHandler := seriesSetup()

	route.Get("/series", seriesHandler.GetAnswerByPosition)
}

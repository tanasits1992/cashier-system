package handlers

import (
	"cashier/app/models"
	"cashier/app/services"
	"cashier/utils"
	"errors"
	"strconv"

	"github.com/gofiber/fiber"
)

type SeriesHandler struct {
	service services.SeriesService
}

func NewSeriesHandler(service services.SeriesService) *SeriesHandler {
	return &SeriesHandler{
		service: service,
	}
}

func (h *SeriesHandler) GetAnswerByPosition(c *fiber.Ctx) {
	queryParam := new(models.SeriesQueryParams)
	if err := c.QueryParser(queryParam); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	if len(queryParam.Positions) == 0 {
		utils.BadRequestResponse(c, "positions is require")
		return
	}

	position := []int{}
	for _, val := range queryParam.Positions {
		p, err := strconv.Atoi(val)

		if err != nil {
			utils.BadRequestResponse(c, err.Error())
			return
		}

		if p < 0 {
			utils.BadRequestResponse(c, errors.New("found position is less than 0").Error())
			return
		}

		position = append(position, p)
	}

	ans := h.service.Calculate(position)

	utils.SuccessResponse(c, ans)
}

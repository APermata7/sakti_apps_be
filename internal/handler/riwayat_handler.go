package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/usecase"
)

type RiwayatHandler struct {
	RiwayatUsecase *usecase.RiwayatUsecase
}

func NewRiwayatHandler(riwayatUsecase *usecase.RiwayatUsecase) *RiwayatHandler {
	return &RiwayatHandler{RiwayatUsecase: riwayatUsecase}
}

func (h *RiwayatHandler) GetRiwayat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	resp, err := h.RiwayatUsecase.GetRiwayat(c.Context(), userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    resp,
	})
}
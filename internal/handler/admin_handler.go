package handler

import (
	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/usecase"
)

type AdminHandler struct {
	AdminUsecase *usecase.AdminUsecase
}

func NewAdminHandler(adminUsecase *usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{AdminUsecase: adminUsecase}
}

func (h *AdminHandler) GetDashboard(c *fiber.Ctx) error {
	stats, err := h.AdminUsecase.GetDashboardStats(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}
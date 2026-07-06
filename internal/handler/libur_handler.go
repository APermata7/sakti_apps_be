package handler

import (
	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/usecase"
)

type LiburHandler struct {
	LiburUsecase *usecase.LiburUsecase
}

func NewLiburHandler(liburUsecase *usecase.LiburUsecase) *LiburHandler {
	return &LiburHandler{LiburUsecase: liburUsecase}
}

func (h *LiburHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateLiburRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	libur, err := h.LiburUsecase.Create(c.Context(), req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    libur,
	})
}

func (h *LiburHandler) GetAll(c *fiber.Ctx) error {
	tahun := c.Query("tahun")
	items, err := h.LiburUsecase.GetAll(c.Context(), tahun)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    items,
	})
}

func (h *LiburHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	libur, err := h.LiburUsecase.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	if libur == nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Libur tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    libur,
	})
}

func (h *LiburHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req domain.UpdateLiburRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	libur, err := h.LiburUsecase.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    libur,
	})
}

func (h *LiburHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.LiburUsecase.Delete(c.Context(), id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Libur berhasil dihapus",
	})
}

func (h *LiburHandler) IsHoliday(c *fiber.Ctx) error {
	tanggal := c.Query("tanggal")
	if tanggal == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Tanggal wajib diisi",
		})
	}

	isHoliday, err := h.LiburUsecase.IsHoliday(c.Context(), tanggal)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"is_holiday": isHoliday,
		},
	})
}
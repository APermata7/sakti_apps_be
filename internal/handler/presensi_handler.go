package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/usecase"
)

type PresensiHandler struct {
	PresensiUsecase *usecase.PresensiUsecase
}

func NewPresensiHandler(presensiUsecase *usecase.PresensiUsecase) *PresensiHandler {
	return &PresensiHandler{PresensiUsecase: presensiUsecase}
}

func (h *PresensiHandler) CheckIn(c *fiber.Ctx) error {
	karyawanID := c.Locals("user_id").(string)

	var req domain.CheckInRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.SelfieURL == "" || req.Latitude == 0 || req.Longitude == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Selfie URL, latitude, dan longitude wajib diisi",
		})
	}

	resp, err := h.PresensiUsecase.CheckIn(c.Context(), karyawanID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	message := ""
	if resp.Status == "terlambat" {
		message = "Anda terlambat, silakan isi alasan"
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    resp,
		"message": message,
	})
}

func (h *PresensiHandler) CheckOut(c *fiber.Ctx) error {
	karyawanID := c.Locals("user_id").(string)

	var req domain.CheckOutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.Latitude == 0 || req.Longitude == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Latitude dan longitude wajib diisi",
		})
	}

	resp, err := h.PresensiUsecase.CheckOut(c.Context(), karyawanID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	message := ""
	if resp.Lembur {
		message = "Lembur " + strconv.FormatFloat(resp.JamLembur, 'f', 1, 64) + " jam tercatat"
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    resp,
		"message": message,
	})
}

func (h *PresensiHandler) GetToday(c *fiber.Ctx) error {
	karyawanID := c.Locals("user_id").(string)

	resp, err := h.PresensiUsecase.GetToday(c.Context(), karyawanID)
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

func (h *PresensiHandler) GetHistory(c *fiber.Ctx) error {
	karyawanID := c.Locals("user_id").(string)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.Query("limit", "30"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	if limit <= 0 {
		limit = 30
	}
	if page <= 0 {
		page = 1
	}

	items, total, err := h.PresensiUsecase.GetHistory(c.Context(), karyawanID, startDate, endDate, status, limit, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"items": items,
			"meta": fiber.Map{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
		},
	})
}

func (h *PresensiHandler) UpdateAlasanTerlambat(c *fiber.Ctx) error {
	karyawanID := c.Locals("user_id").(string)

	var req struct {
		AlasanTerlambat string `json:"alasan_terlambat"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.AlasanTerlambat == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Alasan terlambat wajib diisi",
		})
	}

	err := h.PresensiUsecase.UpdateAlasanTerlambat(c.Context(), karyawanID, req.AlasanTerlambat)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alasan terlambat berhasil disimpan",
	})
}
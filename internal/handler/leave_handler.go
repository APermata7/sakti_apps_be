package handler

import (
	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/usecase"
)

type LeaveHandler struct {
	LeaveUsecase *usecase.LeaveUsecase
}

func NewLeaveHandler(leaveUsecase *usecase.LeaveUsecase) *LeaveHandler {
	return &LeaveHandler{LeaveUsecase: leaveUsecase}
}

func (h *LeaveHandler) CreateLeave(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req struct {
		TipePengajuan  string `json:"tipe_pengajuan"`
		TanggalMulai   string `json:"tanggal_mulai"`
		TanggalSelesai string `json:"tanggal_selesai"`
		Alasan         string `json:"alasan"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.TipePengajuan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Tipe pengajuan wajib diisi (cuti/dispen)",
		})
	}

	if req.TanggalMulai == "" || req.TanggalSelesai == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Tanggal mulai dan tanggal selesai wajib diisi",
		})
	}

	if req.Alasan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Alasan wajib diisi",
		})
	}

	err := h.LeaveUsecase.CreateLeave(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Pengajuan cuti berhasil dikirim",
	})
}

func (h *LeaveHandler) GetStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	leaves, err := h.LeaveUsecase.GetStatus(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    leaves,
	})
}

func (h *LeaveHandler) DownloadSuratCuti(c *fiber.Ctx) error {
	leaveID := c.Params("id")
	if leaveID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID pengajuan wajib diisi",
		})
	}

	pdfBytes, filename, err := h.LeaveUsecase.DownloadSuratCuti(c.Context(), leaveID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return c.Send(pdfBytes)
}

func (h *LeaveHandler) CancelLeave(c *fiber.Ctx) error {
	leaveID := c.Params("id")
	if leaveID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID pengajuan wajib diisi",
		})
	}

	userID := c.Locals("user_id").(string)

	err := h.LeaveUsecase.CancelLeave(c.Context(), leaveID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pengajuan cuti berhasil dibatalkan",
	})
}

func (h *LeaveHandler) ApproveLeave(c *fiber.Ctx) error {
	leaveID := c.Params("id")
	if leaveID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID pengajuan wajib diisi",
		})
	}

	userID := c.Locals("user_id").(string)

	err := h.LeaveUsecase.ApproveLeave(c.Context(), leaveID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pengajuan cuti berhasil disetujui",
	})
}

func (h *LeaveHandler) RejectLeave(c *fiber.Ctx) error {
	leaveID := c.Params("id")
	if leaveID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID pengajuan wajib diisi",
		})
	}

	userID := c.Locals("user_id").(string)

	var req struct {
		Alasan string `json:"alasan"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.Alasan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Alasan penolakan wajib diisi",
		})
	}

	err := h.LeaveUsecase.RejectLeave(c.Context(), leaveID, userID, req.Alasan)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pengajuan cuti ditolak",
	})
}

func (h *LeaveHandler) FinalizeLeave(c *fiber.Ctx) error {
	leaveID := c.Params("id")
	if leaveID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID pengajuan wajib diisi",
		})
	}

	userID := c.Locals("user_id").(string)

	err := h.LeaveUsecase.FinalizeLeave(c.Context(), leaveID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pengajuan cuti berhasil difinalisasi",
	})
}
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/domain"
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

	var req domain.CreateCutiRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.SubTipe == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Sub tipe wajib diisi",
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

	cuti, err := h.LeaveUsecase.CreateLeave(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Pengajuan cuti berhasil",
		"data":    cuti,
	})
}

func (h *LeaveHandler) GetStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	status := c.Query("status")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	items, total, err := h.LeaveUsecase.GetStatus(c.Context(), userID, status, limit, page)
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

func (h *LeaveHandler) GetBalance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	yearParam := c.Query("year")
	var year int
	var err error

	if yearParam != "" {
		year, err = strconv.Atoi(yearParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Format tahun tidak valid",
			})
		}
	} else {
		year = 0
	}

	balance, err := h.LeaveUsecase.GetBalance(c.Context(), userID, year)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    balance,
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

	cuti, err := h.LeaveUsecase.CancelLeave(c.Context(), leaveID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pengajuan cuti berhasil dibatalkan",
		"data":    cuti,
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

	cuti, err := h.LeaveUsecase.ApproveLeave(c.Context(), leaveID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pengajuan cuti berhasil disetujui",
		"data":    cuti,
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

	var req domain.RejectCutiRequest
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

	cuti, err := h.LeaveUsecase.RejectLeave(c.Context(), leaveID, userID, req.Alasan)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pengajuan cuti ditolak",
		"data":    cuti,
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

	cuti, err := h.LeaveUsecase.FinalizeLeave(c.Context(), leaveID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pengajuan cuti berhasil difinalisasi",
		"data":    cuti,
	})
}

func (h *LeaveHandler) GetAllLeaves(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role != "atasan" && role != "hrd" && role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Akses ditolak",
		})
	}

	var req domain.LeaveFilterRequest
	req.Status = c.Query("status")
	req.SubTipe = c.Query("sub_tipe")
	req.StartDate = c.Query("start_date")
	req.EndDate = c.Query("end_date")
	req.Limit, _ = strconv.Atoi(c.Query("limit", "10"))
	req.Page, _ = strconv.Atoi(c.Query("page", "1"))

	items, total, err := h.LeaveUsecase.GetAllLeaves(c.Context(), userID, role, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	totalPages := (total + req.Limit - 1) / req.Limit

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"items": items,
			"meta": fiber.Map{
				"total":       total,
				"page":        req.Page,
				"limit":       req.Limit,
				"total_pages": totalPages,
			},
		},
	})
}

func (h *LeaveHandler) GetApprovalList(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role != "atasan" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Hanya atasan yang dapat mengakses",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	items, total, err := h.LeaveUsecase.GetApprovalList(c.Context(), userID, limit, page)
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

func (h *LeaveHandler) GetFinalizationList(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "hrd" && role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Akses ditolak. Hanya HRD dan Admin yang dapat mengakses",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	items, total, err := h.LeaveUsecase.GetFinalizationList(c.Context(), limit, page)
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
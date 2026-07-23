package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/domain"
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

func (h *AdminHandler) CreateKaryawan(c *fiber.Ctx) error {
	var req domain.CreateKaryawanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid: " + err.Error(),
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email wajib diisi",
		})
	}

	if req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password wajib diisi",
		})
	}

	if len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password minimal 8 karakter",
		})
	}

	if req.NamaLengkap == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama lengkap wajib diisi",
		})
	}

	if req.Role == "" {
		req.Role = "karyawan"
	}

	validRoles := map[string]bool{
		"admin": true, "atasan": true, "hrd": true, "karyawan": true,
	}
	if !validRoles[req.Role] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Role tidak valid. Pilih: admin, atasan, hrd, karyawan",
		})
	}

	if req.LevelJabatan != nil && *req.LevelJabatan != "" {
		validLevels := map[string]bool{
			"staff": true, "officer": true, "spv": true,
			"ka_unit": true, "manager": true, "gm": true, "hrd": true,
		}
		if !validLevels[*req.LevelJabatan] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Level jabatan tidak valid. Pilih: staff, officer, spv, ka_unit, manager, gm, hrd",
			})
		}
	}

	karyawan, err := h.AdminUsecase.CreateKaryawan(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Karyawan berhasil dibuat",
		"data":    karyawan,
	})
}

func (h *AdminHandler) GetAllKaryawan(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search")
	role := c.Query("role")
	status := c.Query("status")

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	karyawanList, total, err := h.AdminUsecase.GetAllKaryawan(c.Context(), page, limit, search, role, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"success": true,
		"data":    karyawanList,
		"meta": fiber.Map{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	})
}

func (h *AdminHandler) GetKaryawan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID karyawan wajib diisi",
		})
	}

	karyawan, err := h.AdminUsecase.GetKaryawanByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if karyawan == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Karyawan tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    karyawan,
	})
}

func (h *AdminHandler) UpdateKaryawan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID karyawan wajib diisi",
		})
	}

	var req domain.UpdateKaryawanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid: " + err.Error(),
		})
	}

	if req.Role != nil && *req.Role != "" {
		validRoles := map[string]bool{
			"admin": true, "atasan": true, "hrd": true, "karyawan": true,
		}
		if !validRoles[*req.Role] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Role tidak valid. Pilih: admin, atasan, hrd, karyawan",
			})
		}
	}

	if req.LevelJabatan != nil && *req.LevelJabatan != "" {
		validLevels := map[string]bool{
			"staff": true, "officer": true, "spv": true,
			"ka_unit": true, "manager": true, "gm": true, "hrd": true,
		}
		if !validLevels[*req.LevelJabatan] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Level jabatan tidak valid. Pilih: staff, officer, spv, ka_unit, manager, gm, hrd",
			})
		}
	}

	if req.StatusKaryawan != nil && *req.StatusKaryawan != "" {
		validStatus := map[string]bool{
			"aktif": true, "nonaktif": true,
		}
		if !validStatus[*req.StatusKaryawan] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Status tidak valid. Pilih: aktif, nonaktif",
			})
		}
	}

	karyawan, err := h.AdminUsecase.UpdateKaryawan(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Karyawan berhasil diupdate",
		"data":    karyawan,
	})
}

func (h *AdminHandler) DeleteKaryawan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID karyawan wajib diisi",
		})
	}

	userID := c.Locals("user_id").(string)
	if id == userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak dapat menghapus akun sendiri",
		})
	}

	err := h.AdminUsecase.DeleteKaryawan(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Karyawan berhasil dinonaktifkan",
	})
}

func (h *AdminHandler) GetPresensiReport(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	items, total, err := h.AdminUsecase.GetPresensiReport(c.Context(), startDate, endDate, status, limit, offset)
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

func (h *AdminHandler) ExportPresensiCSV(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")

	csvData, err := h.AdminUsecase.ExportPresensiCSV(c.Context(), startDate, endDate, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	filename := "laporan_presensi_" + time.Now().Format("20060102_150405") + ".csv"
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename="+filename)
	return c.Send(csvData)
}

func (h *AdminHandler) GetCutiReport(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	items, total, err := h.AdminUsecase.GetCutiReport(c.Context(), startDate, endDate, status, limit, offset)
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

func (h *AdminHandler) ExportCutiCSV(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")

	csvData, err := h.AdminUsecase.ExportCutiCSV(c.Context(), startDate, endDate, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	filename := "laporan_cuti_" + time.Now().Format("20060102_150405") + ".csv"
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename="+filename)
	return c.Send(csvData)
}
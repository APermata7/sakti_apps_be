package handler

import (
	"strconv"
	"time"

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	libur, err := h.LiburUsecase.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    libur,
		"message": "Hari libur berhasil ditambahkan",
	})
}

func (h *LiburHandler) GetAll(c *fiber.Ctx) error {
	tahunStr := c.Query("tahun")
	jenis := c.Query("jenis")
	sumber := c.Query("sumber")
	aktifStr := c.Query("aktif")
	limitStr := c.Query("limit", "30")
	pageStr := c.Query("page", "1")

	var tahun int
	if tahunStr != "" {
		t, err := strconv.Atoi(tahunStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Tahun tidak valid",
			})
		}
		tahun = t
	}

	var aktif *bool
	if aktifStr != "" {
		b := aktifStr == "true"
		aktif = &b
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 30
	}

	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	items, total, err := h.LiburUsecase.GetAll(c.Context(), tahun, jenis, sumber, aktif, limit, page)
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

func (h *LiburHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	libur, err := h.LiburUsecase.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	if libur == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Hari libur tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    libur,
	})
}

func (h *LiburHandler) GetByTahun(c *fiber.Ctx) error {
	tahun, err := strconv.Atoi(c.Params("tahun"))
	if err != nil || tahun <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Tahun tidak valid",
		})
	}

	items, err := h.LiburUsecase.GetByTahun(c.Context(), tahun)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"items": items,
			"meta": fiber.Map{
				"total": len(items),
				"tahun": tahun,
			},
		},
	})
}

func (h *LiburHandler) GetByBulan(c *fiber.Ctx) error {
	bulan := c.Params("bulan")
	if len(bulan) != 7 || bulan[4] != '-' {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format bulan harus YYYY-MM",
		})
	}

	items, err := h.LiburUsecase.GetByBulan(c.Context(), bulan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"items": items,
			"meta": fiber.Map{
				"total": len(items),
				"bulan": bulan,
			},
		},
	})
}

func (h *LiburHandler) CheckTanggal(c *fiber.Ctx) error {
	tanggal := c.Params("tanggal")
	if _, err := time.Parse("2006-01-02", tanggal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format tanggal harus YYYY-MM-DD",
		})
	}

	libur, isLibur, err := h.LiburUsecase.CheckTanggal(c.Context(), tanggal)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if !isLibur || libur == nil {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"is_libur": false,
				"tanggal":  tanggal,
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"is_libur": true,
			"tanggal":  libur.Tanggal.Format("2006-01-02"),
			"nama":     libur.Nama,
			"jenis":    libur.Jenis,
			"aktif":    libur.Aktif,
			"sumber":   libur.Sumber,
		},
	})
}

func (h *LiburHandler) GetAktif(c *fiber.Ctx) error {
	items, err := h.LiburUsecase.GetAktif(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"items": items,
			"meta": fiber.Map{
				"total": len(items),
			},
		},
	})
}

func (h *LiburHandler) GetActiveList(c *fiber.Ctx) error {
	tahunStr := c.Query("tahun")
	bulan := c.Query("bulan")
	jenis := c.Query("jenis")

	var tahun int
	if tahunStr != "" {
		t, err := strconv.Atoi(tahunStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Tahun tidak valid",
			})
		}
		tahun = t
	}

	var items []domain.Libur
	var err error

	if tahun > 0 && bulan != "" {
		items, err = h.LiburUsecase.GetByBulan(c.Context(), bulan)
	} else if tahun > 0 {
		items, err = h.LiburUsecase.GetByTahun(c.Context(), tahun)
	} else if bulan != "" {
		items, err = h.LiburUsecase.GetByBulan(c.Context(), bulan)
	} else {
		items, err = h.LiburUsecase.GetAktif(c.Context())
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if jenis != "" {
		filtered := []domain.Libur{}
		for _, item := range items {
			if item.Jenis == jenis {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    items,
		"meta": fiber.Map{
			"total": len(items),
		},
	})
}

func (h *LiburHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req domain.UpdateLiburRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	libur, err := h.LiburUsecase.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    libur,
		"message": "Hari libur berhasil diupdate",
	})
}

func (h *LiburHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.LiburUsecase.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Hari libur berhasil dihapus",
	})
}

func (h *LiburHandler) Toggle(c *fiber.Ctx) error {
	id := c.Params("id")
	aktif, err := h.LiburUsecase.Toggle(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	status := "dinonaktifkan"
	if aktif {
		status = "diaktifkan"
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Hari libur berhasil " + status,
		"data": fiber.Map{
			"id":    id,
			"aktif": aktif,
		},
	})
}

func (h *LiburHandler) IsHoliday(c *fiber.Ctx) error {
	tanggal := c.Query("tanggal")
	if tanggal == "" {
		tanggal = time.Now().Format("2006-01-02")
	}

	if _, err := time.Parse("2006-01-02", tanggal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format tanggal harus YYYY-MM-DD",
		})
	}

	libur, isLibur, err := h.LiburUsecase.CheckTanggal(c.Context(), tanggal)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if !isLibur || libur == nil {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"is_holiday": false,
				"tanggal":    tanggal,
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"is_holiday": true,
			"tanggal":    tanggal,
			"nama":       libur.Nama,
			"jenis":      libur.Jenis,
		},
	})
}
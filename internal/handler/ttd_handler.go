package handler

import (
	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/usecase"
	"sakti_apps_be/internal/utils"
)

type TTDHandler struct {
	TTDUsecase *usecase.TTDUsecase
}

func NewTTDHandler(ttdUsecase *usecase.TTDUsecase) *TTDHandler {
	return &TTDHandler{TTDUsecase: ttdUsecase}
}

func (h *TTDHandler) UploadTTD(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	karyawanID := c.FormValue("karyawan_id")

	if !(role == "admin" && karyawanID != "") {
		karyawanID = userID
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File image wajib diupload",
		})
	}

	if file.Size > 2*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Ukuran file maksimal 2MB",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membaca file",
		})
	}
	defer src.Close()

	url, err := utils.UploadTTD(src, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal upload file: " + err.Error(),
		})
	}

	ttd, err := h.TTDUsecase.Create(c.Context(), karyawanID, domain.CreateTTDRequest{
		URLTandaTangan: url,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal simpan ke database: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Tanda tangan berhasil diunggah",
		"data":    ttd,
	})
}

func (h *TTDHandler) GetTTD(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	karyawanID := c.FormValue("karyawan_id")

	if !(role == "admin" && karyawanID != "") {
		karyawanID = userID
	}

	ttd, err := h.TTDUsecase.GetByKaryawanID(c.Context(), karyawanID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    ttd,
	})
}

func (h *TTDHandler) UpdateTTD(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	var req struct {
		KaryawanID     string `json:"karyawan_id"`
		URLTandaTangan string `json:"url_tanda_tangan"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.URLTandaTangan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "URL tanda tangan wajib diisi",
		})
	}

	var karyawanID string
	if role == "admin" && req.KaryawanID != "" {
		karyawanID = req.KaryawanID
	} else {
		karyawanID = userID
	}

	ttd, err := h.TTDUsecase.Update(c.Context(), karyawanID, domain.CreateTTDRequest{
		URLTandaTangan: req.URLTandaTangan,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Tanda tangan berhasil diupdate",
		"data":    ttd,
	})
}

func (h *TTDHandler) DeleteTTD(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	karyawanID := c.FormValue("karyawan_id")

	if !(role == "admin" && karyawanID != "") {
		karyawanID = userID
	}

	err := h.TTDUsecase.Delete(c.Context(), karyawanID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Tanda tangan berhasil dihapus",
	})
}

func (h *TTDHandler) VerifyTTD(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	karyawanID := c.FormValue("karyawan_id")

	if !(role == "admin" && karyawanID != "") {
		karyawanID = userID
	}

	result, err := h.TTDUsecase.Verify(c.Context(), karyawanID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}
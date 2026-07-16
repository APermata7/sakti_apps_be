package handler

import (
    "github.com/gofiber/fiber/v2"

    "sakti_apps_be/internal/domain"
    "sakti_apps_be/internal/usecase"
)

type TTDHandler struct {
    TTDUsecase *usecase.TTDUsecase
}

func NewTTDHandler(ttdUsecase *usecase.TTDUsecase) *TTDHandler {
    return &TTDHandler{TTDUsecase: ttdUsecase}
}

func (h *TTDHandler) UploadTTD(c *fiber.Ctx) error {
    karyawanID := c.Locals("user_id").(string)

    var req domain.CreateTTDRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Request tidak valid",
        })
    }

    ttd, err := h.TTDUsecase.Create(c.Context(), karyawanID, req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "Tanda tangan berhasil diunggah",
        "data":    ttd,
    })
}

func (h *TTDHandler) GetTTD(c *fiber.Ctx) error {
    karyawanID := c.Locals("user_id").(string)

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
    karyawanID := c.Locals("user_id").(string)

    var req domain.CreateTTDRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Request tidak valid",
        })
    }

    ttd, err := h.TTDUsecase.Update(c.Context(), karyawanID, req)
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
    karyawanID := c.Locals("user_id").(string)

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
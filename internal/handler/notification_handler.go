package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/usecase"
)

type NotificationHandler struct {
	NotifUsecase *usecase.NotificationUsecase
}

func NewNotificationHandler(notifUsecase *usecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{NotifUsecase: notifUsecase}
}

func (h *NotificationHandler) GetNotifikasi(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if limit > 50 {
		limit = 50
	}

	items, total, err := h.NotifUsecase.GetNotifikasi(c.Context(), userID, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
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

func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	count, err := h.NotifUsecase.GetUnreadCount(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"unread_count": count,
		},
	})
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	id := c.Params("id")

	err := h.NotifUsecase.MarkAsRead(c.Context(), id, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notifikasi ditandai telah dibaca",
	})
}

func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	err := h.NotifUsecase.MarkAllAsRead(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Semua notifikasi ditandai telah dibaca",
	})
}

type RegisterFCMTokenRequest struct {
	FCMToken   string `json:"fcm_token"`
	DeviceID   string `json:"device_id"`
	DeviceType string `json:"device_type"`
}

func (h *NotificationHandler) RegisterFCMToken(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req RegisterFCMTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.FCMToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "FCM token wajib diisi",
		})
	}

	if req.DeviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Device ID wajib diisi",
		})
	}

	err := h.NotifUsecase.RegisterFCMToken(c.Context(), userID, req.FCMToken, req.DeviceID, req.DeviceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "FCM token berhasil disimpan",
	})
}

func (h *NotificationHandler) DeactivateFCMToken(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req struct {
		FCMToken string `json:"fcm_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.FCMToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "FCM token wajib diisi",
		})
	}

	err := h.NotifUsecase.DeactivateFCMToken(c.Context(), userID, req.FCMToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "FCM token berhasil dinonaktifkan",
	})
}
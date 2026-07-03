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
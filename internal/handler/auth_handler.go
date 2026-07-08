package handler

import (
	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/middleware"
	"sakti_apps_be/internal/usecase"
)

type AuthHandler struct {
	AuthUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{AuthUsecase: authUsecase}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email dan password wajib diisi",
		})
	}

	remaining := middleware.GetRemainingAttempts(req.Email)

	resp, err := h.AuthUsecase.Login(c.Context(), req)

	if err != nil {
		middleware.RecordLoginAttempt(req.Email, false)

		remaining = middleware.GetRemainingAttempts(req.Email)

		if remaining == 0 {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Anda telah mencapai batas maksimal percobaan. Akun terkunci selama 5 menit.",
				"remaining_attempts": 0,
				"can_reset": true,
			})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Email atau password salah",
			"remaining_attempts": remaining,
		})
	}

	middleware.RecordLoginAttempt(req.Email, true)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    resp,
	})
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	karyawan, err := h.AuthUsecase.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    karyawan,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logout berhasil",
	})
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req struct {
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.NewPassword == "" || req.ConfirmPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password baru dan konfirmasi wajib diisi",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password minimal 8 karakter",
		})
	}

	if req.NewPassword != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Konfirmasi password tidak sesuai",
		})
	}

	err := h.AuthUsecase.ChangePassword(c.Context(), userID, req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password berhasil diubah",
	})
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email wajib diisi",
		})
	}

	err := h.AuthUsecase.ForgotPassword(c.Context(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Link reset password telah dikirim ke email Anda",
	})
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token           string `json:"token"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.Token == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Token, password baru, dan konfirmasi wajib diisi",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password minimal 8 karakter",
		})
	}

	if req.NewPassword != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Konfirmasi password tidak sesuai",
		})
	}

	err := h.AuthUsecase.ResetPassword(c.Context(), req.Token, req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password berhasil direset",
	})
}
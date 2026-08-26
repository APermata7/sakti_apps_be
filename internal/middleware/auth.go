package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/utils"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak ditemukan",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Format token tidak valid. Gunakan 'Bearer <token>'",
			})
		}

		tokenString := parts[1]

		claims, err := utils.ParseSupabaseJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid: " + err.Error(),
			})
		}

		dbPool, ok := c.Locals("db").(*pgxpool.Pool)
		if !ok || dbPool == nil {
			log.Println("AuthMiddleware: database connection not found")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Internal server error",
			})
		}

		karyawanRepo := repository.NewKaryawanRepo(dbPool)
		karyawan, err := karyawanRepo.GetByEmail(c.Context(), claims.Email)
		if err != nil || karyawan == nil {
			log.Printf("AuthMiddleware: karyawan not found for email: %s", claims.Email)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Data karyawan tidak ditemukan",
			})
		}

		if karyawan.StatusKaryawan != "aktif" {
			log.Printf("AuthMiddleware: akun nonaktif untuk email: %s, status: %s", claims.Email, karyawan.StatusKaryawan)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"code":    "ACCOUNT_INACTIVE",
				"message": "Akun Anda sudah tidak aktif. Silakan hubungi administrator.",
			})
		}

		c.Locals("auth_user_id", claims.Subject)
		c.Locals("user_id", karyawan.ID)
		c.Locals("email", claims.Email)
		c.Locals("token", tokenString)
		c.Locals("karyawan", karyawan)

		var role string

		if claims.AppMetadata != nil {
			if r, ok := claims.AppMetadata["role"].(string); ok {
				role = r
			}
		}

		if role == "" && claims.UserMetadata != nil {
			if r, ok := claims.UserMetadata["role"].(string); ok {
				role = r
			}
		}

		if role == "" {
			role = karyawan.Role
		}

		if role == "" {
			log.Printf("Role tidak ditemukan di metadata untuk user: %s", claims.Email)
			log.Printf("   AppMetadata: %+v", claims.AppMetadata)
			log.Printf("   UserMetadata: %+v", claims.UserMetadata)
		} else {
			log.Printf("Role ditemukan: %s untuk user: %s, karyawan_id: %s", role, claims.Email, karyawan.ID)
		}

		c.Locals("role", role)

		return c.Next()
	}
}

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Role tidak ditemukan",
			})
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak memiliki akses",
		})
	}
}
package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var JWTSecret []byte

func init() {
	_ = godotenv.Load()

	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("Variabel lingkungan JWT_SECRET_KEY tidak diatur. Aplikasi tidak dapat berjalan.")
	}
	JWTSecret = []byte(secret)
}

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Header otentikasi tidak ditemukan"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Format token tidak valid"})
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Metode signing tidak diharapkan")
			}
			return JWTSecret, nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid atau kedaluwarsa"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// --- PERUBAHAN UTAMA DI SINI ---
			// Ubah "user_id" menjadi "userId" agar cocok dengan isi token
			userIdFloat, ok := claims["userId"].(float64)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Claim userId tidak valid dalam token"})
			}

			c.Locals("userId", int64(userIdFloat)) // Nama di Locals boleh tetap snake_case
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Claim token tidak valid"})
	}
}

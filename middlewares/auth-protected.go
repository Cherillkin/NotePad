package middlewares

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Cherillkin/Notepad/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthProtected(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == " " {
			log.Warnf("empty authorization token")

			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warnf("invalid authorization token")

			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}

		tokenStr := tokenParts[1]
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("unexpted signing methos: %v", token.Header["alg"])
			}

			return secret, nil
		})

		if err != nil || !token.Valid {
			log.Warnf("invalid token")

			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}

		userId := token.Claims.(jwt.MapClaims)["id"]

		if err := db.Model(models.User{}).Where("id = ?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warnf("user not found in the db")

			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}

		c.Locals("userId", userId)

		return c.Next()
	}
}

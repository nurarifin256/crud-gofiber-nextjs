package middleware

import (
	"fmt"
	"simple-crud/configs"
	"simple-crud/helpers"
	"simple-crud/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const invalidTokenMessage = "Invalid token"

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return unauthorizedResponse(c)
		}

		tokenString := authHeader[len("Bearer "):]
		if tokenString == "" {
			return unauthorizedResponse(c)
		}

		db := configs.DB.Db

		if !isTokenValid(db, tokenString) {
			return unauthorizedResponse(c)
		}

		token, err := parseToken(tokenString)
		if err != nil {
			return unauthorizedResponse(c)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if !isUserValid(db, int(claims["user_id"].(float64)), c) {
				return unauthorizedResponse(c)
			}
			return c.Next()
		}

		return unauthorizedResponse(c)
	}
}

func unauthorizedResponse(c *fiber.Ctx) error {
	return helpers.ResponseJson(c, fiber.StatusUnauthorized, "error", invalidTokenMessage, []interface{}{})
}

func isTokenValid(db *gorm.DB, tokenString string) bool {
	var userToken models.UserToken
	return db.Where("token = ?", tokenString).First(&userToken).Error == nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
}

func isUserValid(db *gorm.DB, userID int, c *fiber.Ctx) bool {
	var user models.User
	if db.Where("id = ?", userID).First(&user).Error != nil {
		return false
	}
	c.Locals("user", user)
	return true
}

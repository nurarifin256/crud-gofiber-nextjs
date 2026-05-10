package middleware

import (
	"strings"

	users "simple-crud/internal/users"
	"simple-crud/pkg/ctxactor"
	"simple-crud/pkg/response"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	errMisingAuthorization = "Missing Authorization header!"
	errorInvalidHeader     = "Invalid Authorization header format!"
	errorInvalid           = "Invalid or expired token!"
)

type AuthMiddlewareHandler struct {
	service users.Service
	secret  string
}

func NewAuthMiddleware(service users.Service, secret string) *AuthMiddlewareHandler {
	return &AuthMiddlewareHandler{service: service, secret: secret}
}

// Middleware utama
func (h *AuthMiddlewareHandler) Auth(c *fiber.Ctx) error {
	jwtSecret := h.secret
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		zap.S().Error("Missing Authorization Header")
		return response.UnauthorizedErr(errMisingAuthorization)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		zap.S().Errorf("Missing Bearer in Authorization Header: %v", authHeader)
		return response.UnauthorizedErr(errorInvalidHeader)
	}

	tokenString := parts[1]

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			zap.S().Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fiber.ErrUnauthorized
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		zap.S().Errorf("JWT parse error: %v", err)
		return response.UnauthorizedErr(errorInvalid)
	}

	// take nik from claims
	nikUser, ok := claims["nik"].(string)
	if !ok {
		zap.S().Errorf("JWT claims error: %v", ok)
		return response.UnauthorizedErr(errorInvalid)
	}

	// Find user in db
	user, err := h.service.FindUserByNik(c.UserContext(), nikUser)
	if err != nil {
		zap.S().Errorf("User not found: %v", err)
		return response.UnauthorizedErr(errorInvalid)
	}

	// Match token with remember_token (handling pointers)
	if user.RememberToken == nil || *user.RememberToken != tokenString {
		zap.S().Errorf("Token mismatch for user %v", nikUser)
		return response.UnauthorizedErr(errorInvalid)
	}

	// UserID table m_user int64 not pointer
	userID := user.ID
	userNik := user.NIK
	userName := user.Name

	a := ctxactor.Actor{UserID: userID, NIK: userNik, Name: userName}
	c.Locals("actor", a)
	c.SetUserContext(ctxactor.WithActor(c.UserContext(), a))

	return c.Next()
}

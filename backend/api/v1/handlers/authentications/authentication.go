package handler

import (
	users "simple-crud/internal/users"
	"simple-crud/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service users.Service
}

func NewAuthHandler(service users.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

/*
GET /user/find-by-nik/:nik
*/
func (h *AuthHandler) FindUserByNik(c *fiber.Ctx) error {
	nik := c.Params("nik")
	if nik == "" {
		return response.Error(c, fiber.StatusBadRequest, "NIK is required")
	}

	user, err := h.service.FindUserByNik(c.UserContext(), nik)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "User found", user)
}

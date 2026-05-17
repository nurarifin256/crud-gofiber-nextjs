package handler

import (
	"simple-crud/helpers"
	users "simple-crud/internal/users"
	"simple-crud/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service users.Service
}

func NewUserHandler(service users.Service) *UserHandler {
	return &UserHandler{service: service}
}

/*
POST /user/submit
*/
func (h *UserHandler) SubmitUser(c *fiber.Ctx) error {
	var req users.SubmitUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	if err := helpers.V.Struct(req); err != nil {
		return err
	}

	data, err := h.service.SubmitUser(c.Context(), req, "")
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "Success submit user", data)
}

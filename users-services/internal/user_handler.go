package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi11/microservices/users-services/pkg"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req UserRegister
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, pkg.ErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, pkg.ErrorResponse("Validation failed", err.Error()))
		return
	}

	user, err := h.service.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pkg.ErrorResponse("Failed to create user", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, pkg.SuccessResponse("User created successfully", user))
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi11/microservices/gateway/internal/client"
	"github.com/wafi11/microservices/users-services/proto"
)

type UserHandler struct {
	userClient *client.UserClient
}

func NewUserHandler(userClient *client.UserClient) *UserHandler {
	return &UserHandler{userClient: userClient}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req proto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userClient.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

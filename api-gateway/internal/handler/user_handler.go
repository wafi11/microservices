package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wafi11/microservices/api-gateway/internal/client"
	"github.com/wafi11/microservices/api-gateway/pkg"
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
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, pkg.ErrorResponse("Invalid request body", err.Error()))
		return
	}

	user, err := h.userClient.RegisterUser(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pkg.ErrorResponse("Failed to create user", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, pkg.SuccessResponse("User created successfully", user))
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req proto.LoginRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, pkg.ErrorResponse("Invalid request body", err.Error()))
		return
	}

	token, err := h.userClient.LoginUser(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pkg.ErrorResponse("Failed to login user", err.Error()))
		return
	}

	pkg.SetTokenToCookie(c, "access_token", token.Token, "localhost")

	c.JSON(http.StatusCreated, pkg.SuccessResponse("Login successfully", nil))
}

func (h *UserHandler) FindMe(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	// Cast dari any ke string
	userIdStr, ok := userId.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "invalid userId format",
		})
		return
	}
	userIdInt, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid userId",
		})
		return
	}

	resp, err := h.userClient.FindMe(c, &proto.FindMeRequest{
		UserId: int32(userIdInt),
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

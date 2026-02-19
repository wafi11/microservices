package server

import (
	"github.com/gin-gonic/gin"
	"github.com/wafi11/microservices/api-gateway/internal/client"
	"github.com/wafi11/microservices/api-gateway/internal/handler"
	"github.com/wafi11/microservices/api-gateway/pkg"
)

func Routes(r *gin.Engine) {
	userClient, _ := client.NewUserClient("localhost:50051")
	userHandler := handler.NewUserHandler(userClient)

	api := r.Group("/api")
	{
		api.POST("/users", userHandler.CreateUser)
		api.POST("/users/login", userHandler.LoginUser)
		users := api.Use(pkg.AuthMiddleware())
		users.GET("/users/me", userHandler.FindMe)
	}
}

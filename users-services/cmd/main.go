package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wafi11/microservices/users-services/internal"
)

func main() {
	router := gin.New()
	router.Use(gin.Recovery())

	repo := internal.NewUserRepository()
	service := internal.NewUserService(repo)
	handler := internal.NewUserHandler(service)

	router.POST("/api/users", handler.CreateUser)
	s := &http.Server{
		Addr:           ":5001",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Server Listening on port 5001")
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}

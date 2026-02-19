package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wafi11/microservices/api-gateway/server"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:     false,
		AllowOrigins:        []string{"http://localhost:3000"},
		AllowMethods:        []string{"POST", "GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowPrivateNetwork: false,
		AllowHeaders:        []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials:    true,
		MaxAge:              0,
	}))

	server.Routes(r)

	log.Println("HTTP running on :5000")
	if err := r.Run(":5000"); err != nil {
		log.Fatal(err)
	}
	log.Println("HTTP running on :5000")

}

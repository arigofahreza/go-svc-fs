package main

import (
	"go-svc-fs/configs"
	"go-svc-fs/controllers"
	"go-svc-fs/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(CORSMiddleware())

	db, err := configs.SupabaseConfig()
	if err != nil {
		panic(err)
	}
	profileController := controllers.InitProfileController(db)

	mainGroup := router.Group("/api/v1")
	router.Static("/images", "./tmp")
	routes.ProfileRouters(mainGroup, profileController)
	router.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

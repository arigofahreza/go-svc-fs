package main

import (
	"go-svc-fs/configs"
	"go-svc-fs/controllers"
	docs "go-svc-fs/docs"
	"go-svc-fs/routes"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title FS Pictures Service API
// @version 1.0
// @description API for file sytem in profile pictures system
func main() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(CORSMiddleware())
	docs.SwaggerInfo.BasePath = "/"

	db, err := configs.SupabaseConfig()
	if err != nil {
		panic(err)
	}
	profileController := controllers.InitProfileController(db)

	mainGroup := router.Group("/api/v1")
	router.Static("/images", "./tmp")
	routes.ProfileRouters(mainGroup, profileController)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
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

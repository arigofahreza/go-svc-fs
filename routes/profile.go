package routes

import (
	"go-svc-fs/controllers"

	"github.com/gin-gonic/gin"
)

func ProfileRouters(router *gin.RouterGroup, controllers *controllers.ProfileController) {
	group := router.Group("/profile")
	group.GET("image", controllers.GetProfilePicturesControllers)
	group.POST("image", controllers.AddProfilePicturesControllers)
	group.DELETE("image", controllers.DeleteProfilePicturesControllers)
}

package controllers

import (
	"go-svc-fs/models"
	"go-svc-fs/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type ProfileController struct {
	DB *supabase.Client
}

func InitProfileController(db *supabase.Client) *ProfileController {
	return &ProfileController{
		DB: db,
	}
}

func (controller *ProfileController) AddProfilePicturesControllers(c *gin.Context) {
	username := c.PostForm("username")
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "upload image error",
		})
	}
	profileServices := services.InitProfileService(controller.DB, c)
	profileServices.AddProfilePictures(username, file, header)
}

func (controller *ProfileController) GetProfilePicturesControllers(c *gin.Context) {
	username := c.Query("username")
	profileServices := services.InitProfileService(controller.DB, c)
	profileServices.GetProfilePictures(username)
}

func (controller *ProfileController) DeleteProfilePicturesControllers(c *gin.Context) {
	body := models.DeleteImageModel{}
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    400,
			"message": "parsing body error",
		})
	}
	profileServices := services.InitProfileService(controller.DB, c)
	profileServices.DeleteProfilePictures(body.Username, body.Filename)
}

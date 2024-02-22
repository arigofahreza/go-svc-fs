package services

import (
	"encoding/json"
	"fmt"
	"go-svc-fs/utils"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type ProfileService struct {
	DB *supabase.Client
	C  *gin.Context
}

func InitProfileService(db *supabase.Client, c *gin.Context) *ProfileService {
	return &ProfileService{
		DB: db,
		C:  c,
	}
}

func (service *ProfileService) AddProfilePictures(username string, file multipart.File, header *multipart.FileHeader) {
	if username == "" {
		service.C.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "username empty",
		})
		return
	}
	filename := header.Filename
	stringId := fmt.Sprintf("%s_%s", username, filename)
	hashStringId := utils.Md5Hash(stringId)
	out, err := os.Create("./tmp/" + hashStringId + ".png")
	if err != nil {
		service.C.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "creating file error",
		})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		service.C.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "saving file error",
		})
		return
	}

	data := map[string]interface{}{
		"id":         hashStringId,
		"username":   username,
		"created_at": time.Now().Format("2006-01-02 05:32:03"),
	}

	_, _, err = service.DB.From("images").Insert(data, true, "", "", "").Execute()
	if err != nil {
		service.C.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "insert data error",
			"data":    make([]string, 0),
		})
		return
	}

	service.C.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "upload done",
		"data":    data,
	})
}

func (service *ProfileService) GetProfilePictures(username string) {
	if username == "" {
		service.C.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "username empty",
		})
		return
	}
	res, _, err := service.DB.From("images").Select("id", "exact", false).Filter("username", "eq", username).Execute()
	if err != nil {
		service.C.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "get data error",
			"data":    make([]string, 0),
		})
		return
	}
	var jsonDatas []map[string]interface{}
	err = json.Unmarshal(res, &jsonDatas)
	if err != nil {
		service.C.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    make([]string, 0),
		})
		return
	}
	files, err := ioutil.ReadDir("./tmp")
	if err != nil {
		service.C.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "read directory error",
			"data":    make([]string, 0),
		})
		return
	}
	var imageUrls []string
	for _, data := range jsonDatas {
		filename := fmt.Sprintf("%s.png", data["id"])
		for _, file := range files {
			if file.Name() == filename {
				imageUrls = append(imageUrls, fmt.Sprintf("/images/%s", filename))
			}
		}
	}

	service.C.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "get data done",
		"data":    imageUrls,
	})

}

func (service *ProfileService) DeleteProfilePictures(username string, filename string) {
	if username == "" || filename == "" {
		service.C.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "username or filename empty",
		})
		return
	}
	splitName := strings.Split(filename, ".")
	_, _, err := service.DB.From("images").Delete("", "").Filter("id", "eq", splitName[0]).Execute()
	if err != nil {
		service.C.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "delete data error",
			"data":    make([]string, 0),
		})
		return
	}

	data := map[string]interface{}{
		"id":       splitName[0],
		"username": username,
	}

	filepath := fmt.Sprintf("./tmp/%s.png", splitName[0])

	err = os.Remove(filepath)
	if err != nil {
		service.C.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "delete file error",
			"data":    make([]string, 0),
		})
		return
	}

	service.C.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "delete done",
		"data":    data,
	})

}

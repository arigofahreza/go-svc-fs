package models

type DeleteImageModel struct {
	Username string `json:"username" example:"name"`
	Filename string `json:"filename" example:"images.png"`
}

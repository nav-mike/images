package entity

type UploadImageDTO struct {
	Filename string
	File     string // base64 encoded image
	UserId   string
}

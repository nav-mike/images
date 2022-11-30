package entity

type UploadImageDTO struct {
	Filename string
	File     string // Base64 Encoded String
	UserId   string
}

package entity

// UploadImageDTO represents input data (JSON) for saving image
type UploadImageDTO struct {
	Filename string // Filename represents name of the file from the client
	File     string // File represents base64 encoded string. For this project I used https://www.base64encoder.io/image-to-base64-converter/
	UserId   string // UserId represents user id. Used for creating directory for user. In real life it should get from auth middleware (e.g. JWT)
}

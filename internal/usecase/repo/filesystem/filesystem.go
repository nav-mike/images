package filesystem

type FileSystem struct {
	Path string
}

type UploadedImageResponse map[string]string

func NewFileSystem(path string) *FileSystem {
	return &FileSystem{Path: path}
}

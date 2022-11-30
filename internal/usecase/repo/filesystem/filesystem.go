package filesystem

type FileSystem struct {
	Path string
}

func NewFileSystem(path string) *FileSystem {
	return &FileSystem{Path: path}
}

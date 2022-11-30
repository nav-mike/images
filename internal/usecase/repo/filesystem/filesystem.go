package filesystem

// FileSystem represents service to save image to the local filesystem
type FileSystem struct {
	Path string
}

type ValidationError struct {
	Message string
}

func NewFileSystem(path string) *FileSystem {
	return &FileSystem{Path: path}
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

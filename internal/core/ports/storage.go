package ports

import "mime/multipart"

type ImageStorageService interface {
	UploadImage(file *multipart.FileHeader) (string, error)
	DeleteImage(imageURL string) error
}

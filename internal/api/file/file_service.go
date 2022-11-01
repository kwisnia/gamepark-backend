package file

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/uploader"
	"mime/multipart"
)

func UploadImage(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", fmt.Errorf("invalid file")
	}
	filePath, err := uploader.UploadFile("gamepark-images", *file)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

package uploader

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/awsconf"
	"mime/multipart"
	"strings"
)

func UploadFile(bucket string, file multipart.FileHeader) (string, error) {
	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	awsUploader := s3manager.NewUploader(awsconf.AWSSession)
	// check if file is an image
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return "", fmt.Errorf("invalid file type")
	}
	fmt.Println("file type", file.Header.Get("Content-Type"))
	fmt.Println("file size", file.Size)
	if file.Size > 10000000 {
		return "", fmt.Errorf("provided file is too big")
	}
	// get file extension
	fileExtension := strings.Split(file.Filename, ".")[1]
	uploadResult, err := awsUploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(uuid.NewString() + "." + fileExtension),
		ACL:         aws.String("public-read"),
		Body:        fileContent,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}
	return uploadResult.Location, nil
}

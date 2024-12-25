package lib

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

var (
	IMAGE_PATH = os.Getenv("IMAGE_PATH")
)

func FileUpload(file multipart.File, fileHeader multipart.FileHeader) (string, error) {
	basePath := "cmd/web"

	sanitizedFileName := filepath.Base(fileHeader.Filename)
	if err := os.MkdirAll(basePath+IMAGE_PATH, os.ModePerm); err != nil {
		return "", err
	}

	imgPath := strings.Join([]string{IMAGE_PATH, sanitizedFileName}, "")

	imgOut, err := os.Create(basePath + imgPath)
	if err != nil {
		return "", err
	}
	defer imgOut.Close()

	_, err = io.Copy(imgOut, file)
	if err != nil {
		return "", err
	}

	return imgPath, nil
}

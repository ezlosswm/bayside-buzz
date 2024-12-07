package lib

import (
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func FileUpload(path string, file multipart.File, fileHeader multipart.FileHeader) (string, error) {
	// create file
	// orgDir := "/assets/images/organizers/"
	orgDir := strings.Join([]string{"/assets/images/", path,  "/"}, "")
	basePath := "cmd/web"
	// baseDir :="../../cmd/web"

	sanitizedFileName := filepath.Base(fileHeader.Filename)
	if err := os.MkdirAll(basePath+orgDir, os.ModePerm); err != nil {
		slog.Error("Failed to create directories", err)
		return "", err
	}

	imgPath := strings.Join([]string{orgDir, sanitizedFileName}, "")

	imgOut, err := os.Create(basePath + imgPath)
	if err != nil {
		slog.Error("err with path too", err)
		return "", err
	}
	defer imgOut.Close()

	_, err = io.Copy(imgOut, file)
	if err != nil {
		slog.Error("err here too", err)
		return "", nil
	}

	return imgPath, nil
}
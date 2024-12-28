package lib

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/appwrite/sdk-for-go/file"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/storage"
	_ "github.com/joho/godotenv/autoload"
)

func Uploader(s *storage.Storage, imgFile multipart.File, fileHeader multipart.FileHeader) (string, error) {
	sanitizedFileName := filepath.Base(fileHeader.Filename)

	// Create temp file
	tmp, err := os.CreateTemp("", sanitizedFileName)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())

	// Copy the uploaded file to the temp file
	if _, err := io.Copy(tmp, imgFile); err != nil {
		slog.Error("error copying file", "error", err)
		return "", err
	}

	// upload file to appwrite
	resp, err := s.CreateFile(
		os.Getenv("BUCKET_ID"),
		id.Unique(),
		file.NewInputFile(tmp.Name(), sanitizedFileName),
		s.WithCreateFilePermissions([]string{}),
	)
	if err != nil {
		slog.Error("error creating file", "error", err)
		return "", err
	}

	// Close the file
	if err := tmp.Close(); err != nil {
		slog.Error("error closing file", "error", err)
		return "", err
	}

	url := fmt.Sprintf("https://cloud.appwrite.io/v1/storage/buckets/%s/files/%s/view?project=%s", os.Getenv("BUCKET_ID"), resp.Id, os.Getenv("PROJECT_KEY"))

	return url, nil
}

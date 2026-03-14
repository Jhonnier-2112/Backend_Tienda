package storage

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"tienda-backend/internal/core/ports"
)

type LocalStorage struct {
	basePath string
	baseURL  string
}

func NewLocalStorage(basePath, baseURL string) ports.ImageStorageService {
	// Ensure the directory exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.MkdirAll(basePath, os.ModePerm)
	}
	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

func (s *LocalStorage) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", errors.New("no file provided")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}
	defer src.Close()

	// Generate a unique filename using timestamp
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(fileHeader.Filename)
	uniqueFilename := fmt.Sprintf("%d%s", timestamp, ext)

	// Define full save path
	dstPath := filepath.Join(s.basePath, uniqueFilename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("could not create file destination: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("could not save file: %w", err)
	}

	// Return public URL based on the baseURL matching the static Gin route
	imageURL := fmt.Sprintf("%s/%s", s.baseURL, uniqueFilename)
	return imageURL, nil
}

func (s *LocalStorage) DeleteImage(imageURL string) error {
	// Example URL: http://localhost:8080/uploads/1234.jpg
	// We extract just the filename
	fileName := filepath.Base(imageURL)
	filePath := filepath.Join(s.basePath, fileName)

	// Check if file exists before deleting
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File not found, considered successful
	}

	return os.Remove(filePath)
}

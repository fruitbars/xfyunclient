package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ReadImageFile reads an image file and returns its format and binary data.
func ReadImageFile(filePath string) (string, []byte, error) {
	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	// Read the entire file into a byte slice
	imageData, err := io.ReadAll(file)
	if err != nil {
		return "", nil, err
	}

	// Get the file extension and convert it to lowercase
	ext := strings.ToLower(filepath.Ext(filePath))
	var format string

	// Determine the format based on the file extension
	switch ext {
	case ".jpg", ".jpeg":
		format = "jpeg"
	case ".png":
		format = "png"
	case ".gif":
		format = "gif"
	default:
		return "", nil, fmt.Errorf("unsupported image format: %s", ext)
	}

	return format, imageData, nil
}

func ReadImageFileAndReturnBase64(filePath string) (string, string, error) {

	format, imageData, err := ReadImageFile(filePath)
	if err != nil {
		return "", "", err
	}

	// Encode the image data to base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)

	return format, base64Data, nil
}

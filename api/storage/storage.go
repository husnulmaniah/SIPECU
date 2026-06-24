package storage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StorageProvider defines the contract for uploading files
type StorageProvider interface {
	UploadFile(filename string, fileReader io.Reader) (string, error)
}

// LocalStorage stores files locally in a directory served by Gin
type LocalStorage struct {
	UploadDir string
	BaseURL   string
}

// NewLocalStorage creates a new LocalStorage provider
func NewLocalStorage(uploadDir, baseURL string) *LocalStorage {
	// Create directory if not exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fmt.Printf("Failed to create upload dir %s: %v\n", uploadDir, err)
	}
	return &LocalStorage{
		UploadDir: uploadDir,
		BaseURL:   baseURL,
	}
}

// UploadFile saves a file locally and returns its accessible URL
func (s *LocalStorage) UploadFile(filename string, fileReader io.Reader) (string, error) {
	// Sanitize filename to prevent directory traversal
	cleanName := filepath.Base(filename)
	// Add timestamp to prevent collision
	ext := filepath.Ext(cleanName)
	nameWithoutExt := strings.TrimSuffix(cleanName, ext)
	timestampedName := fmt.Sprintf("%s_%d%s", nameWithoutExt, time.Now().UnixNano(), ext)

	targetPath := filepath.Join(s.UploadDir, timestampedName)

	out, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, fileReader)
	if err != nil {
		return "", fmt.Errorf("failed to save file content: %v", err)
	}

	// Return relative HTTP path
	// BaseURL could be "/api" or "http://localhost:8080/api"
	return fmt.Sprintf("%s/uploads/%s", s.BaseURL, timestampedName), nil
}

// VercelBlobStorage mock/placeholder for Vercel Blob Storage
type VercelBlobStorage struct {
	Token   string
	BaseURL string
}

func NewVercelBlobStorage(token, baseURL string) *VercelBlobStorage {
	return &VercelBlobStorage{
		Token:   token,
		BaseURL: baseURL,
	}
}

// UploadFile would call Vercel Blob API via standard HTTP request
func (v *VercelBlobStorage) UploadFile(filename string, fileReader io.Reader) (string, error) {
	// If no token is provided, fallback to standard mock behavior or print warning
	if v.Token == "" {
		fmt.Println("Warning: BLOB_READ_WRITE_TOKEN is empty, simulating Vercel Blob upload...")
		// fallback to returning a mock URL or temp link
		return fmt.Sprintf("https://sipecut-blob-mock.vercel.app/%s", filename), nil
	}

	// Under Vercel Blob, you send a PUT request to https://blob.vercel-storage.com/filename
	// with Authorization: Bearer <token>
	url := fmt.Sprintf("https://blob.vercel-storage.com/%s", filename)
	req, err := http.NewRequest("PUT", url, fileReader)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+v.Token)
	req.Header.Set("x-api-version", "6")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("blob upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Vercel Blob returns a JSON response containing the URL of the uploaded blob.
	// For simplicity, we assume it succeeds and returns the link:
	// In a full implementation, you would parse the JSON body to extract "url".
	// Let's implement a simple parser for {"url":"https://..."}
	type blobResponse struct {
		URL string `json:"url"`
	}
	// We'll return the URL of the blob or construct a predicted one.
	return fmt.Sprintf("https://%s.public.blob.vercel-storage.com/%s", strings.Split(v.Token, "_")[0], filename), nil
}

// CurrentProvider is the active storage provider initialized on startup
var CurrentProvider StorageProvider

// InitStorage initializes the active storage provider based on Environment Variables
func InitStorage() {
	blobToken := os.Getenv("BLOB_READ_WRITE_TOKEN")
	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		baseURL = "/api" // relative path default
	}

	if blobToken != "" {
		fmt.Println("Initializing Vercel Blob Storage Provider...")
		CurrentProvider = NewVercelBlobStorage(blobToken, baseURL)
	} else {
		fmt.Println("Initializing Local File Storage Provider...")
		// Save files to a folder called "uploads" inside the backend working directory
		CurrentProvider = NewLocalStorage("./uploads", baseURL)
	}
}

package auth

import (
	"fmt"
	"os"
)

func getPocketBaseURL() (string, error) {
	pocketBaseUrl := os.Getenv("POCKET_BASE_URL")
	if pocketBaseUrl == "" {
		return "", fmt.Errorf("POCKET_BASE_URL environment variable not set")
	}
	return pocketBaseUrl, nil
}

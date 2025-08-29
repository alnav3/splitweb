package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func AuthWithPassword(identity, password string) (*AuthResponse, error) {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/auth-with-password", pocketBaseUrl)

	requestBody, err := json.Marshal(map[string]string{
		"identity": identity,
		"password": password,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("email or password not valid")
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return nil, err
	}

	return &authResponse, nil
}

func ValidateToken(token, id string) error {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/records/%s", pocketBaseUrl, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token or unauthorized access: status %d", resp.StatusCode)
	}

	return nil
}

func GetUserRecord(token, userID string) (*Record, error) {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/records/%s", pocketBaseUrl, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user record: status %d", resp.StatusCode)
	}

	var record Record
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &record, nil
}

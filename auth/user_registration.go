package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterUser(name, email, password, passwordConfirm string) (*Record, error) {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/records", pocketBaseUrl)

	requestBody, err := json.Marshal(map[string]string{
		"name":            name,
		"email":           email,
		"password":        password,
		"passwordConfirm": passwordConfirm,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response Record
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

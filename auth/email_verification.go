package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RequestEmailChange(newEmail, token string) error {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/request-email-change", pocketBaseUrl)

	requestBody, err := json.Marshal(map[string]string{
		"newEmail": newEmail,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		return fmt.Errorf("There has been an issue while processing your request. Please try again later.")
	}

	return nil
}

func RequestVerificationEmail(email string) error {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/request-verification", pocketBaseUrl)

	requestBody, err := json.Marshal(map[string]string{
		"email": email,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		return fmt.Errorf("There has been an issue while processing your request. Please try again later.")
	}

	return nil
}

func VerifyEmail(token string) error {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/confirm-verification", pocketBaseUrl)

	requestBody, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		return fmt.Errorf("The token allocated is incorrect. Please try again")
	}

	return nil
}

func ConfirmEmailChange(token, password string) error {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/confirm-email-change", pocketBaseUrl)

	requestBody, err := json.Marshal(map[string]string{
		"token":    token,
		"password": password,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	} else if resp.StatusCode == 400 {
		return fmt.Errorf("validation_error")
	} else {
		return fmt.Errorf("unexpected error occurred")
	}
}

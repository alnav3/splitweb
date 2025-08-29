package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RequestPasswordReset(email string) error {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/request-password-reset", pocketBaseUrl)

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
		return fmt.Errorf("Email invalid")
	}

	return nil
}

func ConfirmPasswordReset(token, password, passwordConfirm string) error {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/confirm-password-reset", pocketBaseUrl)

	requestBody, err := json.Marshal(map[string]string{
		"token":           token,
		"password":        password,
		"passwordConfirm": passwordConfirm,
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
		return fmt.Errorf("change password request was invalid")
	}

	return nil
}

func ChangePasswordWithOlderOne(oldPassword, password, passwordConfirm, token, userID string) error {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/records/%s", pocketBaseUrl, userID)

	requestBody, err := json.Marshal(map[string]string{
		"oldPassword":     oldPassword,
		"password":        password,
		"passwordConfirm": passwordConfirm,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("change password with older password request failed")
	}

	return nil
}

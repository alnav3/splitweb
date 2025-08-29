package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RequestEmailChange(identity, token string) error {
	pocketBaseUrl, err := getPocketBaseURL()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/collections/Splitweb_users/request-email-change", pocketBaseUrl)

	requestBody, err := json.Marshal(map[string]string{
		"identity": identity,
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

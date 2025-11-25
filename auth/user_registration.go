package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	databaselogic "github.com/alnav3/splitweb/db/database_logic"
	internalRepository "github.com/alnav3/splitweb/db/internal_repository"
)

func RegisterUser(name, email, password, passwordConfirm string, database *databaselogic.Repository) (*Record, error) {
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

	// step 2: create user in database
	if database != nil && database.Queries != nil {
		err = database.Queries.CreateUser(database.Context, internalRepository.CreateUserParams{
			ID:        response.Id,
			Name:      &response.Name,
			Email:     response.Email,
			AvatarUrl: &response.Avatar,
		})
		if err != nil {
			return nil, err
		}
	}

	return &response, nil
}

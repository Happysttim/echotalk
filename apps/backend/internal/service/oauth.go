package service

import (
	"bytes"
	"encoding/json"
	"net/http"

	"echotalk/internal/config"
)

type GoogleOAuthPayload struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
}

func GoogleToken(code string) (map[string]interface{}, error) {
	config := config.Config

	payload := GoogleOAuthPayload{
		GrantType:    "authorization_code",
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		Code:         code,
		RedirectURI:  config.RedirectURL,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(jsonData)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)

	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return tokenResponse, nil
}

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type UserRequest struct {
    Email  string  `json:"email"`
}

func getContact(email string) (*uuid.UUID, error) {
    apiURL := viper.GetString("api_url")
    httpClient := http.Client{}
    // make request
    u := UserRequest{Email: email}
    userData, err := json.Marshal(u)
    if err != nil {
        return nil, err
    }
    userReq, err := http.NewRequest(http.MethodGet, apiURL+"/users", bytes.NewBuffer(userData))
    if err != nil {
        return nil, err
    }
    userReq.Header.Set("Content-Type", "application/json")
    userReq.Header.Set("Authorization", "Bearer "+viper.GetString("access_token"))
    userResp, err := httpClient.Do(userReq)
    if err != nil {
        return nil, err
    }
    user := &UserResponse{}
    userRespData, err := io.ReadAll(userResp.Body)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(userRespData, user)
    if err != nil {
        return nil, err
    }
    return &user.ID, nil
}

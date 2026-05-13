package requests

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

func GetUser(userID string) (*UserResponse, error) {
    apiURL := viper.GetString("api_url")
    httpClient := http.Client{}
    // make request
	var userReq *http.Request
	if _, err := uuid.Parse(userID); err != nil {
		u := UserRequest{Email: userID}
		userData, err := json.Marshal(u)
		if err != nil {
			return nil, err
		}
		userReq, err = http.NewRequest(http.MethodGet, apiURL+"/users", bytes.NewBuffer(userData))
		if err != nil {
			return nil, err
		}
	} else {
		userReq, err = http.NewRequest(http.MethodGet, apiURL+"/users/"+userID, nil)
		if err != nil {
			return nil, err
		}
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
    return user, nil
}

func GetContactID(email string) (*uuid.UUID, error) {
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

func GetContactEmail(uid uuid.UUID) (string, error) {
    apiURL := viper.GetString("api_url")
    httpClient := http.Client{}
    userReq, err := http.NewRequest(http.MethodGet, apiURL+"/users/"+uid.String(), nil)
    if err != nil {
        return "", err
    }
    userReq.Header.Set("Content-Type", "application/json")
    userReq.Header.Set("Authorization", "Bearer "+viper.GetString("access_token"))
    userResp, err := httpClient.Do(userReq)
    if err != nil {
        return "", err
    }
    user := &UserResponse{}
    userRespData, err := io.ReadAll(userResp.Body)
    if err != nil {
        return "", err
    }
    err = json.Unmarshal(userRespData, user)
    if err != nil {
        return "", err
    }
    return user.Email, nil
}


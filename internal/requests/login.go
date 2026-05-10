package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type LoginRequest struct {
    Email     string  `json:"email"`
    Password  string  `json:"password"`
}
type LoginResponse struct {
    RefreshToken  string  `json:"refresh_token"`
    AccessToken   string  `json:"access_token"`
}

func loginWithPassword(email, password string) error {
    apiURL := viper.GetString("api_url")
    login := LoginRequest{Email: email, Password: password}
    data, err := json.Marshal(login)
    if err != nil {
        return err
    }
    // send credentials to server
    resp, err := http.Post(apiURL+"/login", "application/json", bytes.NewReader(data))
    if err != nil {
        return err
    }
    // make sure the login is valid
    if resp.StatusCode != 200 {
        return errors.New(resp.Status)
    }
    // read body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    // unmarshal response
    var user UserResponse
    err = json.Unmarshal(body, &user)
    if err != nil {
        return err
    }
    // update tokens in config file
    viper.Set("name", user.Name)
    viper.Set("email", user.Email)
    viper.Set("refresh_token", user.RefreshToken)
    viper.Set("access_token", user.AccessToken)
    viper.Set("last_refresh", time.Now().Unix())
    viper.WriteConfig()
    return nil
}

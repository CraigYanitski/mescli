package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
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

func loginWithPassword(email, password string) (bool, error) {
    apiURL := viper.GetString("api_url")
    login := LoginRequest{Email: email, Password: password}
    data, err := json.Marshal(login)
    if err != nil {
        log.Println(err)
        return false, err
    }
    // send credentials to server
    resp, err := http.Post(apiURL+"/login", "application/json", bytes.NewReader(data))
    if err != nil {
        log.Println(err)
        return false, err
    }
    // make sure the login is valid
    if resp.StatusCode != 200 {
        log.Println(resp.Status)
        return false, errors.New(resp.Status)
    }
    // read body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return false, err
    }
    // unmarshal response
    var creds LoginResponse
    err = json.Unmarshal(body, &creds)
    if err != nil {
        log.Println(err)
        return false, err
    }
    // update tokens in config file
    viper.Set("refresh_token", creds.RefreshToken)
    viper.Set("access_token", creds.AccessToken)
    viper.Set("last_refresh", time.Now().Unix())
    viper.WriteConfig()
    return true, nil
}

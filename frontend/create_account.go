package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/CraigYanitski/mescli/client"
	crypt "github.com/CraigYanitski/mescli/cryptography"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type CreateRequest struct {
    Email           string  `json:"email"`
    Name            string  `json:"name"`
    Password        string  `json:"password,omitempty"`
    IdentityKey     string  `json:"identity_key"`
    SignedPrekey    string  `json:"signed_key"`
    SignedKey       string  `json:"signed_prekey"`
}
type UserResponse struct {
    ID              uuid.UUID  `json:"id"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
    Email           string     `json:"email"`
    Name            string     `json:"name"`
    HashedPassword  string     `json:"hashed_password,omitempty"`
    IdentityKey     string     `json:"identity_key"`
    SignedPrekey    string     `json:"signed_key"`
    SignedKey       string     `json:"signed_prekey"`
    Initialised     bool       `json:"initialised"`
    RefreshToken    string     `json:"refresh_token"`
    AccessToken     string     `json:"access_token"`
}

func createAccount(email, password string) (bool, error) {
    // get api url
    apiURL := viper.GetString("api_url")
    // initialise client
    c := &client.Client{}
    c.Initialise(false)
    login := CreateRequest{
        Email: email,
        Name:  "Hi, I'm new here",
        Password:     password,
        IdentityKey:  crypt.EncodeECDSAPublicKey(c.IdentityECDSA()),
        SignedPrekey: crypt.EncodeECDHPublicKey(c.SignedPrekey()),
        SignedKey:    hex.EncodeToString(c.SignedKey),
    }
    data, err := json.Marshal(login)
    if err != nil {
        log.Println(err)
        return false, err
    }
    // send credentials to server
    resp, err := http.Post(apiURL+"/users", "application/json", bytes.NewReader(data))
    if err != nil {
        log.Println(err)
        return false, err
    }
    if resp.StatusCode != 201 {
        log.Printf("status %s: invalid user request", resp.Status)
        return false, nil
    }
    // read body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println(err)
        return false, err
    }
    // unmarshal response
    var user UserResponse
    err = json.Unmarshal(body, &user)
    if err != nil {
        log.Println(err)
        return false, err
    }
    // update config
    viper.Set("name",  user.Name)
    viper.Set("email", user.Email)
    viper.Set("refresh_token", user.RefreshToken)
    viper.Set("access_token",  user.AccessToken)
    viper.Set("last_refresh",  time.Now().Unix())
    viper.WriteConfig()
    return true, nil
}

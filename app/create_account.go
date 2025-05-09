package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/CraigYanitski/mescli/internal/client"
	crypt "github.com/CraigYanitski/mescli/internal/cryptography"
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

func createAccount(name, email, password string) error {
    // get api url
    apiURL := viper.GetString("api_url")
    // initialise client
    c := &client.Client{}
    c.Initialise(false)
    // check if name set
    if name == "" {
        name = "Hi, I'm new here"
    }
    // create request json
    login := CreateRequest{
        Email: email,
        Name:  name,
        Password:     password,
        IdentityKey:  crypt.EncodeECDSAPublicKey(c.IdentityECDSA()),
        SignedPrekey: crypt.EncodeECDHPublicKey(c.SignedPrekey()),
        SignedKey:    hex.EncodeToString(c.SignedKey),
    }
    data, err := json.Marshal(login)
    if err != nil {
        return err
    }
    // send credentials to server
    resp, err := http.Post(apiURL+"/users", "application/json", bytes.NewReader(data))
    if err != nil {
        return err
    }
    if resp.StatusCode != 201 {
        return fmt.Errorf("status %s: invalid user request", resp.Status)
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
    // update config
    viper.Set("name",  user.Name)
    viper.Set("email", user.Email)
    viper.Set("refresh_token", user.RefreshToken)
    viper.Set("access_token",  user.AccessToken)
    viper.Set("last_refresh",  time.Now().Unix())
    viper.WriteConfig()
    return nil
}

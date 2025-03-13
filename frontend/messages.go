package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/CraigYanitski/mescli/client"
	"github.com/CraigYanitski/mescli/cryptography"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type UserKeyPacket struct {
    UserID         uuid.UUID  `json:"user_id"`
    IdentityKey    string     `json:"identity_key"`
    SignedPrekey   string     `json:"signed_prekey"`
    SignedKey      string     `json:"signed_key"`
    OnetimePrekey  string     `json:"onetime_prekey"`
}
type MessageRequest struct {
    UserID    uuid.UUID  `json:"user_id"`
    SenderID  uuid.UUID  `json:"sender_id"`
    Message   string     `json:"message"`
}
type MessageResponse struct {
    ID         uuid.UUID  `json:"id"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
    UserID    uuid.UUID  `json:"user_id"`
    SenderID  uuid.UUID  `json:"sender_id"`
    Message   string     `json:"message"`
}

func addContact(email string) error {
    return nil
}

func sendMessage(user uuid.UUID, message string) error {
    apiURL := viper.GetString("api_url")
    httpClient := http.Client{}
    c := client.Client{}
    c.Initialise(false)
    // get user key packet
    u := UserResponse{ID: user}
    userData, err := json.Marshal(u)
    if err != nil {
        return err
    }
    userReq, err := http.NewRequest(http.MethodGet, apiURL+"/users", bytes.NewBuffer(userData))
    userReq.Header.Set("Content-Type", "application/json")
    userReq.Header.Set("Authorization", "Bearer "+viper.GetString("access_token"))
    userResp, err := httpClient.Do(userReq)
    if err != nil {
        return err
    }
    userKeys := &UserKeyPacket{}
    keyData, err := io.ReadAll(userResp.Body)
    if err != nil {
        return err
    }
    err = json.Unmarshal(keyData, userKeys)
    if err != nil {
        return err
    }
    userIK := cryptography.DecodeECDSAPublicKey(userKeys.IdentityKey)
    ikECDH, err := userIK.ECDH()
    if err != nil {
        return err
    }
    encryptedMsg, err := c.SendMessage(message, []string{}, ikECDH, false)
    if err != nil {
        return err
    }
    msg := MessageRequest{
        UserID: user,
        Message: encryptedMsg,
    }
    msgData, err := json.Marshal(msg)
    if err != nil {
        return fmt.Errorf("error marshalling message JSON: %s", err)
    }
    // send request to server
    msgReq, err := http.NewRequest(http.MethodPost, apiURL+"/messages", bytes.NewBuffer(msgData))
    if err != nil {
        return fmt.Errorf("error making message request", err)
    }
    msgReq.Header.Set("Content-Type", "application/json")
    msgReq.Header.Set("Authorization", "Bearer "+viper.GetString("access_token"))
    msgResp, err := httpClient.Do(msgReq)
    defer msgResp.Body.Close()
    // check if request successful
    if msgResp.StatusCode != 201 {
        return errors.New("error: update not successful")
    }
    return nil
}

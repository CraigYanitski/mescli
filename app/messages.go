package main

import (
	"bytes"
	"crypto/ecdsa"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CraigYanitski/mescli/internal/client"
	"github.com/CraigYanitski/mescli/internal/cryptography"
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
    UserID              uuid.UUID  `json:"user_id"`
    SenderID            uuid.UUID  `json:"sender_id"`
    Message             string     `json:"message"`
    SenderIdentityKey   string     `json:"sender_identity_key"`
    SenderEphemeralKey  string     `json:"sender_ephemeral_key"`
}
type MessageResponse struct {
    ID                  uuid.UUID       `json:"id"`
    CreatedAt           time.Time       `json:"created_at"`
    UpdatedAt           time.Time       `json:"updated_at"`
    UserID              uuid.UUID       `json:"user_id"`
    SenderID            uuid.UUID       `json:"sender_id"`
    SenderIdentityKey   sql.NullString  `json:"sender_identity_key"`
    SenderEphemeralKey  sql.NullString  `json:"sender_ephemeral_key"`
    Message             string          `json:"message"`
}

func addContact(email string) (*client.MessagePacketJSON, error) {
    apiURL := viper.GetString("api_url")
    httpClient := http.Client{}
    u := client.Client{}
    u.Initialise(false)
    // get contact key packet
    senderID, err := getContact(email)
    if err != nil {
        return nil, err
    }
    keyReqStruct := UserResponse{ID: *senderID}
    data, err := json.Marshal(keyReqStruct)
    if err != nil {
        return nil, err
    }
    // send request to server
    keyReq, err := http.NewRequest(http.MethodGet, apiURL+"/users/crypto/"+(*senderID).String(), bytes.NewBuffer(data))
    if err != nil {
        return nil, err
    }
    keyReq.Header.Set("Content-Type", "application/json")
    keyReq.Header.Set("Authorization", "Bearer "+viper.GetString("access_token"))
    keyResp, err := httpClient.Do(keyReq)
    if err != nil {
        return nil, err
    }
    senderKeys := &UserKeyPacket{}
    keyRespData, err := io.ReadAll(keyResp.Body)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(keyRespData, senderKeys)
    if err != nil {
        return nil, err
    }
    senderKeyPacket := &client.PrekeyPacketJSON{
        IdentityKey: senderKeys.IdentityKey,
        SignedPrekey: senderKeys.SignedPrekey,
        SignedKey: senderKeys.SignedKey,
        OnetimePrekey: senderKeys.OnetimePrekey,
    }
    messageJSON := u.InitiateX3DH(senderKeyPacket, *senderID, false)
    return messageJSON, nil
}

func getUserIdentityKey(user uuid.UUID) (*ecdsa.PublicKey, error) {
    apiURL := viper.GetString("api_url")
    httpClient := http.Client{}
    // get user key packet
    u := UserResponse{ID: user}
    userData, err := json.Marshal(u)
    if err != nil {
        return nil, err
    }
    userReq, err := http.NewRequest(http.MethodGet, apiURL+"/users/identity/"+user.String(), bytes.NewBuffer(userData))
    if err != nil {
        return nil, err
    }
    userReq.Header.Set("Content-Type", "application/json")
    userReq.Header.Set("Authorization", "Bearer "+viper.GetString("access_token"))
    userResp, err := httpClient.Do(userReq)
    if err != nil {
        return nil, err
    }
    userKeys := &UserKeyPacket{}
    keyData, err := io.ReadAll(userResp.Body)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(keyData, userKeys)
    if err != nil {
        return nil, err
    }
    identityKey := cryptography.DecodeECDSAPublicKey(userKeys.IdentityKey)
    //userIK, err := identityKey.ECDH()
    //if err != nil {
    //    return nil, err
    //}
    return identityKey, nil
}

func sendMessage(contactID uuid.UUID, contactX3DHpacket *client.MessagePacketJSON, message string) error {
    apiURL := viper.GetString("api_url")
    httpClient := http.Client{}
    c := client.Client{}
    c.Initialise(false)
    // get contact identity key (only need identity key for exchanging messages)
    contactIK, err := getUserIdentityKey(contactID)  // TODO: change input to string (and function to decode string)
    if err != nil {
        return fmt.Errorf("error getting contact identity key: %s", err)
    }
    // encrypt message and marshal request JSON
    encryptedMsg, err := c.SendMessage(message, contactIK, contactID, false)
    if err != nil {
        return err
    }
    // check messagePacketJSON
    if contactX3DHpacket == nil {
        contactX3DHpacket = &client.MessagePacketJSON{
            IdentityKey: cryptography.EncodeECDSAPublicKey(contactIK),
            EphemeralKey: "",
        }
    }
    // marshal request JSON
    msg := MessageRequest{
        UserID: contactID,
        Message: encryptedMsg,
        SenderIdentityKey: contactX3DHpacket.IdentityKey,
        SenderEphemeralKey: contactX3DHpacket.EphemeralKey,
    }
    msgData, err := json.Marshal(msg)
    if err != nil {
        return fmt.Errorf("error marshalling message JSON: %s", err)
    }
    // send request to server
    msgReq, err := http.NewRequest(http.MethodPost, apiURL+"/messages", bytes.NewBuffer(msgData))
    if err != nil {
        return fmt.Errorf("error making message request: %s", err)
    }
    msgReq.Header.Set("Content-Type", "application/json")
    msgReq.Header.Set("Authorization", "Bearer "+viper.GetString("access_token"))
    msgResp, err := httpClient.Do(msgReq)
    if err != nil {
        return err
    }
    defer msgResp.Body.Close()
    // check if request successful
    if msgResp.StatusCode != 201 {
        return errors.New("error: update not successful")
    }
    return nil
}

func getMessages() (messages []MessageResponse, err error) {
    messages = []MessageResponse{}
    apiURL := viper.GetString("api_url")
    httpClient := http.Client{}
    c := client.Client{}
    c.Initialise(false)
    // send GET request to server
    req, err := http.NewRequest(http.MethodGet, apiURL+"/messages", nil)
    if err != nil {
        return
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+viper.GetString("access_token"))
    resp, err := httpClient.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        err = errors.New("error: cannot retrieve messages")
        return
    }
    messagesData, err := io.ReadAll(resp.Body)
    if err != nil {
        return
    }
    messagesSlice := &[]MessageResponse{}
    err = json.Unmarshal(messagesData, messagesSlice)
    if err != nil {
        return
    }
    for _, message := range *messagesSlice {
        // get sender conversation from config
        senderMessages := viper.GetStringSlice("contacts."+message.SenderID.String()+".messages")
        senderEncryptedMessages := viper.GetStringSlice("contacts."+message.SenderID.String()+".encrypted_messages")
        senderEncryptedMessages = append(senderEncryptedMessages, message.Message)
        viper.Set("contacts."+message.SenderID.String()+".encrypted_messages", senderEncryptedMessages)
        // check if X3DH initiated
        if message.SenderEphemeralKey.Valid && message.SenderIdentityKey.Valid {
            err = c.CompleteX3DH(
                &client.MessagePacketJSON{
                    IdentityKey: message.SenderIdentityKey.String,
                    EphemeralKey: message.SenderEphemeralKey.String,
                },
                message.SenderID,
                false,
            )
            if err != nil {
                senderMessages = append(senderMessages, message.Message)
                viper.Set("contacts."+message.SenderID.String()+".messages", senderMessages)
                continue
            }
        }
        // get the sender identity key
        senderIK, err := getUserIdentityKey(message.SenderID)
        if err != nil {
            senderMessages = append(senderMessages, message.Message)
            viper.Set("contacts."+message.SenderID.String()+".messages", senderMessages)
            continue
        }
        decryptedMessage, err := c.ReceiveMessage(
            message.Message, 
            senderIK, 
            message.SenderID, 
            false)
        if err != nil {
            senderMessages = append(senderMessages, message.Message)
            viper.Set("contacts."+message.SenderID.String()+".messages", senderMessages)
            continue
        }
        senderMessages = append(senderMessages, decryptedMessage)
        viper.Set("contacts."+message.SenderID.String()+".messages", senderMessages)
        messages = append(messages, message)
    }
    viper.WriteConfig()
    return
}

func writeMessages(messages map[string][]string) bool {
    messageBytes, err := json.MarshalIndent(messages, "", "    ")
    if err != nil {
        log.Println(err)
        return false
    }
    err = os.WriteFile("./.messages", messageBytes, 0644)
    if err != nil {
        log.Println(err)
        return false
    }
    return true
}


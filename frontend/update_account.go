package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	crypt "github.com/CraigYanitski/mescli/cryptography"
	"github.com/spf13/viper"
)

type UpdateRequest struct {
    Email         string  `json:"email"`
    Password      string  `json:"password"`
    Name          string  `json:"name"`
    IdentityKey   string  `json:"identity_key"`
    SignedPrekey  string  `json:"signed_prekey"`
    SignedKey     string  `json:"signed_key"`
}


func updateAccount(name, email, password string) error {
    //get api_url
    apiURL := viper.GetString("api_url")
    // get user_cryptographic keys
    ik := crypt.DecodeECDSAPrivateKey(viper.GetString("identity_key"))
    IK := crypt.EncodeECDSAPublicKey(&ik.PublicKey)
    spk := crypt.DecodeECDHPrivateKey(viper.GetString("signed_prekey"))
    SPK := crypt.EncodeECDHPublicKey(spk.PublicKey())
    SK := viper.GetString("signed_key")
    // create JSON to send as request
    user := UpdateRequest{
        Email: email,
        Password: password,
        Name: name,
        IdentityKey: IK,
        SignedPrekey: SPK,
        SignedKey: SK,
    }
    data, err := json.Marshal(user)
    if err != nil {
        return err
    }
    // send request to server
    client := &http.Client{}
    req, err := http.NewRequest(http.MethodPut, apiURL+"/users", bytes.NewBuffer(data))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+viper.GetString("access_token"))
    resp, err := client.Do(req)
    defer resp.Body.Close()
    // check if request successful
    if resp.StatusCode != 200 {
        return errors.New("error: update not successful")
    }
    // save config changes
    viper.Set("name", name)
    viper.Set("email", email)
    return nil
}


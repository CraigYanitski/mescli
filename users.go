package main

import (
	//"crypto/ecdh"
	//"crypto/ecdsa"
	//"crypto/x509"
	//"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CraigYanitski/mescli/internal/database"
	"github.com/google/uuid"
)

type InitUser struct {
    Email           string  `json:"email"`
    HashedPassword  string  `json:"hashed_password,omitempty"`
    IdentityKey     string  `json:"identity_key"`
    SignedPrekey    string  `json:"signed_key"`
    SignedKey       string  `json:"signed_prekey"`
}
type User struct {
    ID              uuid.UUID  `json:"id"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
    Email           string     `json:"email"`
    HashedPassword  string     `json:"hashed_password,omitempty"`
    IdentityKey     string     `json:"identity_key"`
    SignedPrekey    string     `json:"signed_key"`
    SignedKey       string     `json:"signed_prekey"`
}

type PrekeyPacketJSON struct {
    IdentityKey   string  `json:"identity_key"`
    SignedPrekey  string  `json:"signed_prekey"`
    SignedKey     string  `json:"signed_key"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    u := &InitUser{}
    err := decoder.Decode(u)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "unable to unmarshal user", err)
        return
    }

    //idkBytes, err := x509.MarshalPKIXPublicKey(&u.IdentityKey)
    //if err != nil {
    //    respondWithError(w, http.StatusInternalServerError, "error marshalling identity key", err)
    //    return
    //}
    //spkBytes, err := x509.MarshalPKIXPublicKey(&u.SignedPrekey)
    //if err != nil {
    //    respondWithError(w, http.StatusInternalServerError, "error marshalling signed prekey", err)
    //    return
    //}
    //skBytes, err := x509.MarshalPKIXPublicKey(&u.SignedKey)
    //if err != nil {
    //    respondWithError(w, http.StatusInternalServerError, "error marshalling signed key", err)
    //    return
    //}

    params := database.CreateUserParams{
        Email: u.Email,
        HashedPassword: u.HashedPassword,
        IdentityKey: u.IdentityKey,//hex.EncodeToString(idkBytes),
        SignedPrekey: u.SignedPrekey,//hex.EncodeToString(spkBytes),
        SignedKey: u.SignedKey,//hex.EncodeToString(skBytes),
    }
    createdUser, err := cfg.dbQueries.CreateUser(r.Context(), params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "error adding to database", err)
        return
    }

    createdUser.HashedPassword = ""
    respondWithJSON(w, http.StatusCreated, User(createdUser))
    return
}

func (cfg *apiConfig) handleGetUserKeyPacket(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    u := &User{}
    err := decoder.Decode(u)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "unable to unmarshal user", err)
    }

    userKeys, err := cfg.dbQueries.GetUserKeyPacket(r.Context(), u.ID)
    if err != nil {
        respondWithError(
            w, 
            http.StatusInternalServerError, 
            fmt.Sprintf("error finding %s to database", u.ID), 
            err,
        )
        return
    }

    respondWithJSON(w, http.StatusOK, PrekeyPacketJSON(userKeys))
    return
}


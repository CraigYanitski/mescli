package main

import (
	//"crypto/ecdh"
	//"crypto/ecdsa"
	//"crypto/x509"
	//"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	crypt "github.com/CraigYanitski/mescli/cryptography"
	"github.com/CraigYanitski/mescli/internal/auth"
	"github.com/CraigYanitski/mescli/internal/database"
	"github.com/google/uuid"
)

type InitUser struct {
    Email           string  `json:"email"`
    Name            string  `json:"name"`
    Password        string  `json:"password,omitempty"`
    IdentityKey     string  `json:"identity_key"`
    SignedPrekey    string  `json:"signed_key"`
    SignedKey       string  `json:"signed_prekey"`
}
type User struct {
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
}
type ValidUser struct {
    User
    AccessToken   string  `json:"access_token"`
    RefreshToken  string  `json:"refresh_token"`
}

type PrekeyPacketJSON struct {
    IdentityKey   string  `json:"identity_key"`
    SignedPrekey  string  `json:"signed_prekey"`
    SignedKey     string  `json:"signed_key"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
    // Unmarshal request JSON
    decoder := json.NewDecoder(r.Body)
    u := &InitUser{}
    err := decoder.Decode(u)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "unable to unmarshal user", err)
        return
    }
    // verify required fields are set
    if u.Email == "" {
        respondWithError(w, http.StatusBadRequest, "error: email required to create user", nil)
        return
    } else if u.Password == "" {
        respondWithError(w, http.StatusBadRequest, "error: hashed password required to create user", nil)
        return
    } else if u.IdentityKey == "" {
        respondWithError(w, http.StatusBadRequest, "error: hashed identity key required to create user", nil)
        return
    } else if u.SignedPrekey == "" {
        respondWithError(w, http.StatusBadRequest, "error: signed prekey required to create user", nil)
        return
    } else if u.SignedKey == "" {
        respondWithError(w, http.StatusBadRequest, "error: signed key required to create user", nil)
        return
    }
    // set name if unset
    if u.Name == "" {
        // respondWithError(w, http.StatusBadRequest, "error: name required to create user", nil)
        // return
        u.Name = strings.Split(u.Email, "@")[0]
    }
    //u.IdentityKey = []byte{0, 0}
    //u.SignedPrekey = []byte{0, 0}
    //u.SignedKey = []byte{0, 0}

    hash, err := crypt.HashPassword(u.Password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "failed to hash password", err)
        return
    }

    params := database.CreateUserParams{
        Email: u.Email,
        Name: u.Name,
        HashedPassword: hash,
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

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
    // unmarshal the POST JSON and verify required fields are valid
    decoder := json.NewDecoder(r.Body)
    u := &InitUser{}
    err := decoder.Decode(u)
    if (err != nil) || (u.Email == "") || (u.Password == "") {
        respondWithError(
            w, 
            http.StatusInternalServerError, 
            fmt.Sprintf("error decoding JSON with email '%s' and password '%s'", u.Email, u.Password), 
            err,
        )
        return
    }

    // determine JWT duration
    duration := time.Hour

    // search for user in database using their email
    foundUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), u.Email)
    if (err != nil) || (foundUser.HashedPassword == "") {
        respondWithError(w, http.StatusNotFound, "error finding user", err)
        return
    }

    // make user JWT token
    token, err := auth.MakeJWT(foundUser.ID, cfg.secret, duration)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "error making JWT token", err)
        return
    }

    // generate refresh token
    refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), foundUser.ID)
    if err != nil {
        rtExpiresAt := time.Now().AddDate(0, 0, 60)
        rt, err := auth.MakeRefreshToken()
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "", err)
            return
        }
        params := database.CreateRefreshTokenParams{Token: rt, UserID: foundUser.ID, ExpiresAt: rtExpiresAt}
        refreshToken, err = cfg.dbQueries.CreateRefreshToken(r.Context(), params)
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "error creating refresh token", err)
            return
        }
    }

    // recast database user to validated one, adding JWT
    validUser := &ValidUser{}
    validUser.User = User(foundUser)
    validUser.AccessToken = token
    validUser.RefreshToken = refreshToken.Token

    // check validity of password
    if !crypt.CheckPasswordHash(u.Password, validUser.HashedPassword) {
        respondWithError(w, http.StatusUnauthorized, "password incorrect", err)
        return
    }

    // empty password field to remove from marshalled JSON
    validUser.HashedPassword = ""

    // respond with user JSON
    respondWithJSON(w, http.StatusOK, validUser)
    return
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
    // check user authentication
    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "", err)
        return
    }
    id, err := auth.ValidateJWT(token, cfg.secret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, token, err)
        return
    }

    // unmarshal the POST JSON and verify required fields are valid
    decoder := json.NewDecoder(r.Body)
    u := &InitUser{}
    err = decoder.Decode(u)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "unable to unmarshal request", err)
        return
    } else if u.Email == "" || u.Password == "" {
        respondWithError(
            w, 
            http.StatusInternalServerError, 
            fmt.Sprintf("error decoding JSON with email '%s' and password '%s'", u.Email, u.Password), 
            nil,
        )
        return
    } else if u.IdentityKey == "" {
        respondWithError(w, http.StatusBadRequest, "error: hashed identity key required for user update", nil)
        return
    } else if u.SignedPrekey == "" {
        respondWithError(w, http.StatusBadRequest, "error: signed prekey required for user update", nil)
        return
    } else if u.SignedKey == "" {
        respondWithError(w, http.StatusBadRequest, "error: signed key required for user update", nil)
        return
    }

    hash, err := crypt.HashPassword(u.Password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "failed to hash password", err)
        return
    }

    params := database.UpdateUserParams{
        ID: id,
        Email: u.Email,
        Name: u.Name,
        HashedPassword: hash,
        IdentityKey: u.IdentityKey,//hex.EncodeToString(idkBytes),
        SignedPrekey: u.SignedPrekey,//hex.EncodeToString(spkBytes),
        SignedKey: u.SignedKey,//hex.EncodeToString(skBytes),
    }
    createdUser, err := cfg.dbQueries.UpdateUser(r.Context(), params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "error adding to database", err)
        return
    }

    createdUser.HashedPassword = ""
    respondWithJSON(w, http.StatusCreated, User(createdUser))
    return
}


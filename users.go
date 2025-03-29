package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CraigYanitski/mescli/internal/auth"
	crypt "github.com/CraigYanitski/mescli/internal/cryptography"
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
    OnetimePrekey   string  `json:"onetime_prekey"`
}
type User struct {
    ID              uuid.UUID  `json:"id"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
    Email           string     `json:"email"`
    Name            string     `json:"name"`
    HashedPassword  string     `json:"hashed_password,omitempty"`
    Initialised     bool       `json:"initialised"`
}
type ValidUser struct {
    User
    AccessToken   string  `json:"access_token"`
    RefreshToken  string  `json:"refresh_token"`
}

type CryptoKey struct {
    IdentityKey     string     `json:"identity_key"`
    CreatedAt       time.Time  `json:"created_at,omitempty"`
    UpdatedAt       time.Time  `json:"updated_at,omitempty"`
    UserID          uuid.UUID  `json:"user_id"`
    SignedPrekey    string     `json:"signed_key"`
    SignedKey       string     `json:"signed_prekey"`
    OnetimePrekey   string     `json:"onetime_prekey"`
}

type PrekeyPacketJSON struct {
    IdentityKey    string  `json:"identity_key"`
    SignedPrekey   string  `json:"signed_prekey"`
    SignedKey      string  `json:"signed_key"`
    OnetimePrekey  string  `json:"onetime_prekey"`
}

func (cfg *apiConfig) authenticationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // check user authentication
        token, err := auth.GetBearerToken(r.Header)
        if err != nil {
            respondWithError(w, http.StatusUnauthorized, "", err)
            return
        }

        _, err = auth.ValidateJWT(token, cfg.secret)
        if err != nil {
            respondWithError(w, http.StatusUnauthorized, token, err)
            return
        }

        next.ServeHTTP(w, r)
    })
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

    createParams := database.CreateUserParams{
        Email: u.Email,
        Name: u.Name,
        HashedPassword: hash,
        //IdentityKey: u.IdentityKey,//hex.EncodeToString(idkBytes),
        //SignedPrekey: u.SignedPrekey,//hex.EncodeToString(spkBytes),
        //SignedKey: u.SignedKey,//hex.EncodeToString(skBytes),
    }
    createdUser, err := cfg.dbQueries.CreateUser(r.Context(), createParams)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "error adding to users database", err)
        return
    }

    cryptoParams := database.CreateKeyPacketParams{
        UserID: createdUser.ID,
        IdentityKey: u.IdentityKey,
        SignedPrekey: u.SignedPrekey,
        SignedKey: u.SignedKey,
        OnetimePrekey: u.OnetimePrekey,
    }
    _, err = cfg.dbQueries.CreateKeyPacket(r.Context(), cryptoParams)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "error adding to crypto_keys database", err)
        return
    }

    createdUser.HashedPassword = ""
    respondWithJSON(w, http.StatusCreated, User(createdUser))
    return
}

func (cfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {
    // get user ID from request
    userID, err := uuid.Parse(r.PathValue("userID"))
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "unable to parse user ID", err)
        return
    }
    // check user authentication (will wrap in middleware)
    //token, err := auth.GetBearerToken(r.Header)
    //if err != nil {
    //    respondWithError(w, http.StatusUnauthorized, "", err)
    //    return
    //}
    //_, err = auth.ValidateJWT(token, cfg.secret)
    //if err != nil {
    //    respondWithError(w, http.StatusUnauthorized, token, err)
    //    return
    //}

    user, err := cfg.dbQueries.GetUser(r.Context(), userID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "unable to get user", err)
        return
    }

    // return user
    respondWithJSON(w, http.StatusOK, User(user))
}

func (cfg *apiConfig) handleGetUserByEmail(w http.ResponseWriter, r *http.Request) {
    // check user authentication (will wrap in middleware)
    //token, err := auth.GetBearerToken(r.Header)
    //if err != nil {
    //    respondWithError(w, http.StatusUnauthorized, "", err)
    //    return
    //}
    //_, err = auth.ValidateJWT(token, cfg.secret)
    //if err != nil {
    //    respondWithError(w, http.StatusUnauthorized, token, err)
    //}

    // ensure request contains email
    u := &User{}
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(u)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "unable to unmarshal user", err)
        return
    }
    if u.Email == "" {
        respondWithError(w, http.StatusBadRequest, "need to supply user email", err)
        return
    }

    // get user
    user, err := cfg.dbQueries.GetUserByEmail(r.Context(), u.Email)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "unable to get user by email", err)
        return
    }

    // return user JSON
    respondWithJSON(w, http.StatusOK, User(user))
}

func (cfg *apiConfig) handleGetUserKeyPacket(w http.ResponseWriter, r *http.Request) {
    // get user ID from request
    userID, err := uuid.Parse(r.PathValue("userID"))
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "unable to parse user ID", err)
        return
    }
    // check user authentication (will wrap in middleware)
    //token, err := auth.GetBearerToken(r.Header)
    //if err != nil {
    //    respondWithError(w, http.StatusUnauthorized, "", err)
    //    return
    //}
    //_, err = auth.ValidateJWT(token, cfg.secret)
    //if err != nil {
    //    respondWithError(w, http.StatusUnauthorized, token, err)
    //    return
    //}

    // get user request
    //decoder := json.NewDecoder(r.Body)
    //u := &User{}
    //err = decoder.Decode(u)
    //if err != nil {
    //    respondWithError(w, http.StatusInternalServerError, "unable to unmarshal user", err)
    //    return
    //}
    //if u.Email != "" {
    //    foundUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), u.Email)
    //    if err != nil {
    //        respondWithError(w, http.StatusInternalServerError, "unable to find user in DB by email", err)
    //    } else {
    //        u.ID = foundUser.ID
    //    }
    //}

    // make request for key packet
    userKeyPacket, err := cfg.dbQueries.GetUserKeyPacket(r.Context(), userID)
    if err != nil {
        respondWithError(
            w, 
            http.StatusInternalServerError, 
            fmt.Sprintf("error finding %s to database", userID), 
            err,
        )
        return
    }

    respondWithJSON(w, http.StatusOK, CryptoKey(userKeyPacket))
    return
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
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
    //if (err != nil) || (foundUser.HashedPassword == "") {
    //    respondWithError(w, http.StatusNotFound, "error finding user", err)
    //    return
    //}

    // check validity of password
    if !crypt.CheckPasswordHash(u.Password, foundUser.HashedPassword) {
        respondWithError(w, http.StatusUnauthorized, "password incorrect", err)
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
            respondWithError(w, http.StatusInternalServerError, "error making refresh token", err)
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

    // empty password field to remove from marshalled JSON
    validUser.HashedPassword = ""

    // respond with user JSON
    respondWithJSON(w, http.StatusOK, validUser)
    return
}

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
    // check user authentication (will wrap in middleware)
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
        // IdentityKey: u.IdentityKey,//hex.EncodeToString(idkBytes),
        // SignedPrekey: u.SignedPrekey,//hex.EncodeToString(spkBytes),
        // SignedKey: u.SignedKey,//hex.EncodeToString(skBytes),
    }
    updatedUser, err := cfg.dbQueries.UpdateUser(r.Context(), params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "error updating users database", err)
        return
    }

    cryptoParams := database.UpdateKeyPacketParams{
        UserID: id,
        IdentityKey: u.IdentityKey,
        SignedPrekey: u.SignedPrekey,
        SignedKey: u.SignedKey,
        OnetimePrekey: u.OnetimePrekey,
    }
    _, err = cfg.dbQueries.UpdateKeyPacket(r.Context(), cryptoParams)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "error updating crypto_keys database", err)
        return
    }

    updatedUser.HashedPassword = ""
    respondWithJSON(w, http.StatusCreated, User(updatedUser))
    return
}


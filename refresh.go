package main

import (
	"net/http"
	"time"

	"github.com/CraigYanitski/mescli/internal/auth"
)

type Token struct {
    Token  string  `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "missing token", err)
        return
    }

    refreshToken, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), token)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "invalid entry", err)
        return
    }

    newToken, err := auth.MakeJWT(refreshToken.ID, cfg.secret, time.Hour)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "unauthorised", err)
    }

    respondWithJSON(w, http.StatusOK, Token{newToken})
    return
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "missing token", err)
        return
    }

    err = cfg.dbQueries.RevokeRefreshToken(r.Context(), token)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "invalid entry", err)
        return
    }

    respondWithJSON(w, http.StatusNoContent, nil)
    return
}


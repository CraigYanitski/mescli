package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/CraigYanitski/mescli/internal/auth"
	"github.com/CraigYanitski/mescli/internal/database"
	"github.com/google/uuid"
)

type InitMessage struct {
    UserID    uuid.UUID  `json:"user_id"`
    //SenderID  uuid.UUID  `json:"sender_id"`
    Message   string     `json:"message"`
}
type Message struct {
    ID                  uuid.UUID       `json:"id"`
    CreatedAt           time.Time       `json:"created_at"`
    UpdatedAt           time.Time       `json:"updated_at"`
    UserID              uuid.UUID       `json:"user_id"`
    SenderID            uuid.UUID       `json:"sender_id"`
    SenderIdentityKey   sql.NullString  `json:"sender_identity_key"`
    SenderEphemeralKey  sql.NullString  `json:"sender_ephemeral_key"`
    Message             string          `json:"message"`
}

func (cfg *apiConfig) handleCreateMessage(w http.ResponseWriter, r *http.Request) {
    // check authentication
    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "unauthorised", err)
        return
    }
    id, err := auth.ValidateJWT(token, cfg.secret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, token, err)
        return
    }

    // unmarshal POST JSON
    decoder := json.NewDecoder(r.Body)
    m := &Message{}
    err = decoder.Decode(m)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "error decoding request", err)
        return
    }
    if m.Message == "" {
        respondWithError(w, http.StatusBadRequest, "need message to create entry", err)
        return
    }

    params := database.CreateMessageParams{
        UserID: m.UserID,
        SenderID: id,
        Message: m.Message,
    }
    createdMessage, err := cfg.dbQueries.CreateMessage(r.Context(), params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "error adding to messages database", err)
        return
    }

    respondWithJSON(w, http.StatusCreated, Message(createdMessage))
}

func (cfg *apiConfig) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
    // check authentication
    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "unauthorised", err)
        return
    }
    id, err := auth.ValidateJWT(token, cfg.secret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, token, err)
        return
    }
    
    // get messages
    messages, err := cfg.dbQueries.GetMessages(r.Context(), id)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "error getting messages from database", err)
    }

    // delete messages from database
    newMessages := []Message{}
    for _, message := range messages {
        newMessages = append(newMessages, Message(message))
        _, err = cfg.dbQueries.DeleteMessage(r.Context(), message.ID)
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "error deleting message from server", err)
        }
    }

    respondWithJSON(w, http.StatusOK, newMessages)
}


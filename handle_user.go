package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Suryarpan/rss-agg/internal/auth"
	"github.com/Suryarpan/rss-agg/internal/database"
	"github.com/google/uuid"
)

func (apiCgf *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error while parsing user data")
		return
	}
	user, err := apiCgf.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %s", err))
		return
	}
	convUser := convertDbUser(user)
	respondWithJson(w, 201, convUser)
}

func (apiCgf *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Auth error: %s", err))
		return
	}

	user, err := apiCgf.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %s", err))
		return
	}

	respondWithJson(w, 200, convertDbUser(user))
}

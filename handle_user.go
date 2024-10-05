package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, 200, convertDbUser(user))
}

func (apiCfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get posts: %s", err))
		return
	}
	convPosts := make([]Post, 0, len(posts))
	for _, post := range posts {
		convPosts = append(convPosts, convertDbPost(post))
	}
	respondWithJson(w, 200, convPosts)
}

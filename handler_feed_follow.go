package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Suryarpan/rss-agg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, u database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error while parsing feed follow %s", err))
		return
	}
	feedFollow, err := apiCfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    u.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't follow the feed %s", err))
		return
	}
	respondWithJson(w, 201, convertDbFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), u.ID)
	if err != nil {
		respondWithError(w, 400, "Couldn't get followed feeds of the current user")
		return
	}
	convFeedFollows := make([]FeedFollows, 0, len(feedFollows))
	for _, feedFollow := range feedFollows {
		convFeedFollows = append(convFeedFollows, convertDbFeedFollow(feedFollow))
	}
	respondWithJson(w, 200, convFeedFollows)
}

func (apiCfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	feedFolowIdStr := chi.URLParam(r, "feedFollowId")
	feedFollowId, err := uuid.Parse(feedFolowIdStr)
	if err != nil {
		respondWithError(w, 400, "Error while parsing data")
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: u.ID,
	})
	if err != nil {
		respondWithError(w, 403, "Feed deleted or it does not below to user")
		return
	}
	respondWithJson(w, 200, struct{ message string }{message: "Deleted successfully"})
}

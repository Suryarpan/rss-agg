package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Suryarpan/rss-agg/internal/database"
	"github.com/google/uuid"
)

func (apiCgf *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error while parsing feed data")
		return
	}
	feed, err := apiCgf.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %s", err))
		return
	}
	respondWithJson(w, 201, convertDbFeed((feed)))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Coudn't get feeds: %s", err))
		return
	}
	convFeeds := make([]Feed, 0, len(feeds))

	for _, feed := range feeds {
		convFeeds = append(convFeeds, convertDbFeed(feed))
	}

	respondWithJson(w, 200, convFeeds)
}
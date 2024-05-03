package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Demianeen/chirpy/internal/auth"
	"github.com/Demianeen/chirpy/internal/database"
)

func (config *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	type respose struct {
		database.User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}

	apiKey, err := auth.GetApiToken(r.Header)
	if err != nil || apiKey != config.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate API key")
		return
	}

	user, err := config.db.GetUserById(params.Data.UserID)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}
	user.IsChirpyRed = true

	_, err = config.db.UpdateUserById(user.Id, user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

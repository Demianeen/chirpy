package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Demianeen/chirpy/internal/auth"
	"github.com/Demianeen/chirpy/internal/database"
)

func (config *apiConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email               string `json:"email"`
		Password            string `json:"password"`
		ExpirationInSeconds int    `json:"expires_in_seconds"`
	}
	type respose struct {
		database.User
		JwtToken     string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := config.db.GetUserByEmail(params.Email)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	jwtString, err := auth.GenerateJwt(user.Id, nil, config.jwtSecret)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate jwt")
		return
	}

	refreshToken, err := config.db.CreateRefreshToken(user.Id)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token")
		return
	}

	respondWithJson(w, http.StatusOK, respose{
		User: database.User{
			Id:          user.Id,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		JwtToken:     jwtString,
		RefreshToken: refreshToken.Token,
	})
}

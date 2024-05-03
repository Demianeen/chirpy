package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/Demianeen/chirpy/internal/auth"
	"github.com/Demianeen/chirpy/internal/database"
)

func (config *apiConfig) handleRefreshJwtToken(w http.ResponseWriter, r *http.Request) {
	type respose struct {
		Token string `json:"token"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token")
		return
	}

	oldRefreshToken, err := config.db.GetRefreshTokenData(tokenString)
	if errors.Is(err, database.ErrNotExist) {
		respondWithError(w, http.StatusUnauthorized, "specified token couldn't be found")
		return
	}

	err = config.db.ValidateRefreshToken(oldRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	jwtString, err := auth.GenerateJwt(oldRefreshToken.UserId, nil, config.jwtSecret)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate jwt")
		return
	}

	respondWithJson(w, http.StatusOK, respose{
		Token: jwtString,
	})
}

func (config *apiConfig) handleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	type respose struct {
		Token string `json:"token"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token")
		return
	}

	oldRefreshToken, err := config.db.GetRefreshTokenData(tokenString)
	if errors.Is(err, database.ErrNotExist) {
		respondWithError(w, http.StatusUnauthorized, "specified token couldn't be found")
		return
	}

	err = config.db.ValidateRefreshToken(oldRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	config.db.RevokeRefreshToken(oldRefreshToken.Id)

	respondWithJson(w, http.StatusOK, struct{}{})
}

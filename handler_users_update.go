package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Demianeen/chirpy/internal/auth"
	"github.com/Demianeen/chirpy/internal/database"
)

func (config *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	claims, err := auth.ParseJwt(tokenString, config.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	userId, err := claims.GetSubject()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Invalid JWT subject")
		return
	}

	parsedUserId, err := strconv.Atoi(userId)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Invalid user id")
		return
	}

	user, err := config.db.GetUserById(parsedUserId)
	if err != nil {
		if err == database.ErrNotExist {
			respondWithError(w, http.StatusNotFound, "User not exist")
			return
		}
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	user.Email = params.Email
	user.HashedPassword, err = auth.HashPassword(params.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}
	newUser, err := config.db.UpdateUserById(parsedUserId, user)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJson(w, http.StatusOK, database.GetPublicUser(newUser))
}

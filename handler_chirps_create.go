package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Demianeen/chirpy/internal/auth"
	"github.com/Demianeen/chirpy/internal/database"
)

func replaceProfanity(msg string) string {
	words := strings.Fields(msg)
	bannedWords := map[string]string{
		"kerfuffle": "****",
		"sharbert":  "****",
		"fornax":    "****",
	}

	validatedWords := make([]string, 0, len(words))
	for _, word := range words {
		if replaceValue, ok := bannedWords[strings.ToLower(word)]; ok {
			validatedWords = append(validatedWords, replaceValue)
			continue
		}
		validatedWords = append(validatedWords, word)
	}
	return strings.Join(validatedWords, " ")
}

func (config *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
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

	_, err = config.db.GetUserById(parsedUserId)
	if err != nil {
		if err == database.ErrNotExist {
			respondWithError(w, http.StatusNotFound, "User not exist")
			return
		}
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	if len(params.Body) > 120 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	newCrisp, err := config.db.CreateChirp(params.Body, parsedUserId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respondWithJson(w, http.StatusCreated, newCrisp)
}

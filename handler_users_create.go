package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Demianeen/chirpy/internal/database"
)

func (config *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
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

	newUser, err := config.db.CreateUser(params.Email, params.Password)
	if err != nil {
		if err == database.ErrAlreadyExist {
			respondWithError(w, http.StatusConflict, "User already exist")
			return
		}
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	respondWithJson(w, http.StatusCreated, database.GetPublicUser(newUser))
}

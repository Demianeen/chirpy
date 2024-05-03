package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Demianeen/chirpy/internal/auth"
)

func (config *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
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

	chirpIdString := r.PathValue("chirpId")
	chirpId, err := strconv.Atoi(chirpIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	if chirpId != parsedUserId {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp")
		return
	}

	err = config.db.DeleteChirp(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't delete chirp")
		return
	}

	w.WriteHeader(http.StatusOK)
}

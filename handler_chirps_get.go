package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/Demianeen/chirpy/internal/database"
)

const (
	ascSortOrder  = "asc"
	descSortOrder = "desc"
)

func (config *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue("chirpId")
	chirpId, err := strconv.Atoi(chirpIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := config.db.GetChirpById(chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "The chirp wasn't found")
		return
	}

	respondWithJson(w, http.StatusOK, dbChirp)
}

func (config *apiConfig) handleRetrieveChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := config.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// optional query param to sort chirps by author id
	authorIdString := r.URL.Query().Get("author_id")
	if authorIdString != "" {
		authorId, err := strconv.Atoi(authorIdString)
		if err != nil {
			log.Println(err)
			respondWithError(w, http.StatusInternalServerError, "Invalid user id")
			return
		}
		authorOnlyChirps := make([]database.Chirp, 0)
		for _, chirp := range chirps {
			if chirp.AuthorId == authorId {
				authorOnlyChirps = append(authorOnlyChirps, chirp)
			}
		}
		chirps = authorOnlyChirps
	}

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == descSortOrder {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id > chirps[j].Id
		})
		// asc order in every other case
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		})
	}
	respondWithJson(w, http.StatusOK, chirps)
}

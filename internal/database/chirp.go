package database

type Chirp struct {
	Body     string `json:"body"`
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newChirpId := len(dbData.Chirps) + 1
	newChirp := Chirp{
		Id:       newChirpId,
		AuthorId: authorId,
		Body:     body,
	}
	dbData.Chirps[newChirpId] = newChirp
	db.writeDB(dbData)
	return newChirp, nil
}

func (db *DB) DeleteChirp(chirpId int) error {
	dbData, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbData.Chirps, chirpId)
	db.writeDB(dbData)
	return nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirpById(id int) (Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := dbData.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbData.Chirps))
	for _, chirpValue := range dbData.Chirps {
		chirps = append(chirps, chirpValue)
	}
	return chirps, nil
}

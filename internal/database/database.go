package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	ErrNotExist     = errors.New("resource does not exist")
	ErrAlreadyExist = errors.New("resource already exist")
)

type DB struct {
	mu   *sync.RWMutex
	path string
}

type DBStructure struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RefreshTokens map[int]RefreshToken `json:"refresh_tokens"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	newBD := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	err := newBD.ensureDB()
	if err != nil {
		return nil, err
	}

	return newBD, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = db.writeDB(DBStructure{
				Chirps:        make(map[int]Chirp),
				Users:         make(map[int]User),
				RefreshTokens: make(map[int]RefreshToken),
			})
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbRawData, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbData := DBStructure{}
	err = json.Unmarshal(dbRawData, &dbData)
	if err != nil {
		return DBStructure{}, err
	}
	return dbData, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	json, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, json, 0644)
	return err
}

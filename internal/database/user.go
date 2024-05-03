package database

import (
	"errors"

	"github.com/Demianeen/chirpy/internal/auth"
)

type User struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	Id             int    `json:"id"`
	IsChirpyRed    bool   `json:"is_chirpy_red"`
}

func GetPublicUser(user User) User {
	return User{
		Email: user.Email,
		Id:    user.Id,
	}
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExist
	}

	newUserId := len(dbData.Users) + 1
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return User{}, err
	}

	newUser := User{
		Id:             newUserId,
		Email:          email,
		HashedPassword: hashedPassword,
		IsChirpyRed:    false,
	}
	dbData.Users[newUserId] = newUser
	db.writeDB(dbData)

	return newUser, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}

func (db *DB) GetUserById(id int) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbData.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}
	return user, nil
}

func (db *DB) UpdateUserById(userId int, newUser User) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	dbData.Users[userId] = newUser
	db.writeDB(dbData)

	return newUser, nil
}

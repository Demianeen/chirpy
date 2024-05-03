package database

import (
	"errors"
	"time"

	"github.com/Demianeen/chirpy/internal/auth"
)

type RefreshToken struct {
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Token     string    `json:"token"`
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	IsRevoked bool      `json:"is_revoked"`
}

func (db *DB) CreateRefreshToken(userId int) (RefreshToken, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	refreshTokenString, err := auth.GenerateRefreshToken(userId)
	if err != nil {
		return RefreshToken{}, err
	}

	newId := len(dbData.RefreshTokens) + 1
	refreshToken := RefreshToken{
		Id:        newId,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		UserId:    userId,
		Token:     refreshTokenString,
		IsRevoked: false,
	}

	dbData.RefreshTokens[newId] = refreshToken
	db.writeDB(dbData)
	return refreshToken, nil
}

func (db *DB) RevokeRefreshToken(tokenId int) (RefreshToken, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	token := dbData.RefreshTokens[tokenId]
	token.IsRevoked = true
	dbData.RefreshTokens[tokenId] = token

	db.writeDB(dbData)
	return token, nil
}

func (db *DB) GetRefreshTokenData(tokenString string) (RefreshToken, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	for _, dbToken := range dbData.RefreshTokens {
		if dbToken.Token == tokenString {
			return dbToken, nil
		}
	}

	return RefreshToken{}, ErrNotExist
}

func (db *DB) ValidateRefreshToken(refreshToken RefreshToken) error {
	if refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.IsRevoked {
		return errors.New("refresh token have expired. Please login again")
	}

	return nil
}

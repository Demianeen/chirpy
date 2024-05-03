package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRefreshToken(userId int) (string, error) {
	const size = 32

	bytes := make([]byte, size)
	_, err := rand.Reader.Read(bytes)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)
	return token, nil
}

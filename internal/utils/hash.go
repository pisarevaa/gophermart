package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

func GetPasswordHash(password string, secretKey string) (string, error) {
	str := []byte(password)
	str = append(str, []byte(secretKey)...)
	h := sha256.New()
	_, err := h.Write(str)
	if err != nil {
		return "", err
	}
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return sha, nil
}

func CheckPasswordHash(password string, hash string, secretKey string) (bool, error) {
	hashNew, err := GetPasswordHash(password, secretKey)
	if err != nil {
		return false, err
	}
	if hash != hashNew {
		return false, errors.New("wrong password")
	}
	return true, nil
}

package helper

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

func GenAdminPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password+os.Getenv("PASSWORD_SECRET_ADMIN")), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareAdminPassword(password string, hash string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+os.Getenv("PASSWORD_SECRET_ADMIN")))
	if err != nil {
		return err
	}

	return nil
}

func GenUserPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password+os.Getenv("PASSWORD_SECRET_USER")), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareUserPassword(password string, hash string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+os.Getenv("PASSWORD_SECRET_USER")))
	if err != nil {
		return err
	}

	return nil
}

// pkg/auth: Password hashing with bcrypt.
package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const cost = bcrypt.DefaultCost

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

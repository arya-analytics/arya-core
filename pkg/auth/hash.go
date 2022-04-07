package auth

import "golang.org/x/crypto/bcrypt"

const hashCost = bcrypt.DefaultCost

func GenerateFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(hash), err
}

func compareHashAndPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return Error{
			Type:    ErrorTypeInvalidCredentials,
			Message: "Invalid credentials",
			Base:    err,
		}
	}
	return err
}

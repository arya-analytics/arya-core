package auth

import "golang.org/x/crypto/bcrypt"

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

func GenerateFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

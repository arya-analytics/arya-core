package auth

import "golang.org/x/crypto/bcrypt"

const hashCost = bcrypt.DefaultCost

// HashPassword hashes a password using the bcrypt algorithm.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(hash), err
}

// ValidatePassword checks if password matches hash.
// Returns Error with type ErrorTypeInvalidCredentials if the password does not match the hash.
func ValidatePassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return Error{
			Type:    ErrorTypeInvalidCredentials,
			Message: "Invalid credentials.",
			Base:    err,
		}
	}
	return err
}

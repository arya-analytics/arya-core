package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

// |||| PARAMETERS ||||

func tokenExpiration() time.Duration {
	return 1 * time.Hour
}

func tokenSecret() []byte {
	return []byte("ExtremelySecretKey")
}

// |||| NEW ||||

// NewToken generates a new JWT with the PK of the user as the claim Issuer.
func NewToken(userPK uuid.UUID) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    userPK.String(),
		ExpiresAt: time.Now().Add(tokenExpiration()).Unix(),
	})
	return claims.SignedString(tokenSecret())
}

// |||| PARSE ||||

func parseToken(token string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret(), nil
	})
	return claims, err
}

// |||| VALIDATE ||||

// ValidateToken validates if the token matches the prov
func ValidateToken(token string) error {
	if _, err := parseToken(token); err != nil {
		return Error{
			Type:    ErrorTypeInvalidCredentials,
			Message: "Invalid token",
			Base:    err,
		}
	}
	return nil
}

package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrJWTSecretNotSet      = errors.New("JWT_SECRET environment variable is not set")
	ErrUnexpectedSignMethod = errors.New("unexpected signing method")
	ErrTokenExpired         = errors.New("token expired")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidTokenClaims   = errors.New("invalid token claims")
	ErrInvalidTokenIssuer   = errors.New("invalid token issuer")
)

type TokenClaims struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(ID int, username, email string, expiryHours int) (string, error) {
	jwtSecret, err := getSecretKey()
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(time.Duration(expiryHours) * time.Hour)
	claims := &TokenClaims{
		ID:       ID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func VerifyToken(tokenString string) (*TokenClaims, error) {
	jwtSecret, err := getSecretKey()
	if err != nil {
		return nil, err
	}

	claims := &TokenClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSignMethod
		}
		return jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if err := validateClaims(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func getSecretKey() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, ErrJWTSecretNotSet
	}
	return []byte(secret), nil
}

func validateClaims(claims *TokenClaims) error {
	if claims.ID <= 0 {
		return fmt.Errorf("%w: invalid ID", ErrInvalidTokenClaims)
	}

	if claims.Username == "" {
		return fmt.Errorf("%w: missing username", ErrInvalidTokenClaims)
	}

	if claims.Email == "" {
		return fmt.Errorf("%w: missing email", ErrInvalidTokenClaims)
	}

	if claims.Issuer != os.Getenv("JWT_ISSUER") {
		return fmt.Errorf("%w: issuer mismatch", ErrInvalidTokenIssuer)
	}

	return nil
}

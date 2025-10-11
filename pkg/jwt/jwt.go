package jwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tacheraSasi/go-api-starter/internals/models"
)

type Claims struct {
	User models.User `json:"user"`
	jwt.RegisteredClaims
}

// ValidateToken validates a JWT string and returns claims
func ValidateToken(tokenString string, jwtSecret []byte) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Check signing method (important for security!)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims if token is valid
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateToken generates a new JWT token
func GenerateToken(user models.User, jwtSecret []byte, jwtExpiresIn string) (string, error) {
	expiresIn, err := strconv.Atoi(jwtExpiresIn)
	if err != nil {
		expiresIn = 24
	}
	claims := &Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiresIn))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}


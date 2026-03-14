package utils

import (
	"errors"
	"time"

	"tienda-backend/internal/core/domain"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user *domain.User, duration time.Duration, jwtSecret string) (string, error) {
	expirationTime := time.Now().Add(duration)

	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))

	return tokenString, err
}

func GenerateTokenPair(user *domain.User, jwtSecret string) (*domain.TokenPair, error) {
	// Access Token: 15 minutes
	accessToken, err := GenerateToken(user, 15*time.Minute, jwtSecret)
	if err != nil {
		return nil, err
	}

	// Refresh Token: 7 days
	refreshToken, err := GenerateToken(user, 7*24*time.Hour, jwtSecret)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ValidateToken(tokenString, jwtSecret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

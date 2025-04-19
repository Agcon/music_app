package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTManager interface {
	Generate(userID int64) (string, error)
	Validate(token string) (int64, error)
	GetTTL() time.Duration
}

type jwtManager struct {
	secretKey string
	ttl       time.Duration
}

func NewJWTManager(secretKey string, ttl time.Duration) JWTManager {
	return &jwtManager{secretKey, ttl}
}

func (j *jwtManager) Generate(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(j.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtManager) Validate(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	return userID, nil
}

func (j *jwtManager) GetTTL() time.Duration {
	return j.ttl
}

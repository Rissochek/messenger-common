package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

type TokenClaims struct {
	UserId int
	jwt.StandardClaims
}

type JWTManager struct {
	SecretKey     string
	TokenDuration time.Duration
	logger *log.Logger
}

func NewJWTManager(token_duration time.Duration, secretKey string) *JWTManager {
	return &JWTManager{TokenDuration: token_duration, SecretKey: secretKey}
}

func (manager *JWTManager) VerifyTokenWithClaims(token string) (*TokenClaims, error) {
	tokenJWT, err := jwt.ParseWithClaims(
		token,
		&TokenClaims{},
		func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				manager.logger.Errorf("failed to verify token")
				return nil, errors.New("failed to verify token")
			}

			return []byte(manager.SecretKey), nil
		},
	)
	if err != nil {
		manager.logger.Errorf("invalid token: %v", err)
		return nil, errors.New("invalid token")
	}

	claims, ok := tokenJWT.Claims.(*TokenClaims)
	if !ok {
		manager.logger.Errorf("invalid token claims")
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func ExtractToken(bearerToken string, logger *log.Logger) (string, error) {
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		logger.Errorf("failed to extract token: invalid format")
		return "", errors.New("invalid token format")
	}

	token := strings.TrimPrefix(bearerToken, "Bearer ")
	return token, nil
}
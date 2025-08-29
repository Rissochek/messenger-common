package middleware

import (
	"context"

	"github.com/Rissochek/messenger-common/auth"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Verifier interface {
	VerifyTokenWithClaims(token string) (*auth.TokenClaims, error)
}

type BlacklistChecker interface {
	BlacklistCheck(ctx context.Context, token string) (bool, error)
}

type TokenValidateMiddleware struct {
	Verifier
	BlacklistChecker
	logger *log.Logger
}

func NewTokenValidateMiddleware(verifier Verifier, blacklist BlacklistChecker) *TokenValidateMiddleware {
	return &TokenValidateMiddleware{Verifier: verifier, BlacklistChecker: blacklist}
}

func (middleware *TokenValidateMiddleware) ValidateToken(context *gin.Context) {
	bearerToken := context.GetHeader("Authorization")
	if bearerToken == "" {
		context.JSON(401, gin.H{"error": "authorization header is empty"})
		context.Abort()
		return
	}

	token, err := auth.ExtractToken(bearerToken, middleware.logger)
	if err != nil {
		context.JSON(401, gin.H{"error": "token format is invalid"})
		context.Abort()
		return
	}

	ok, err := middleware.BlacklistChecker.BlacklistCheck(context, token)
	if err != nil {
		context.JSON(500, gin.H{"error": "failed to check token"})
		return
	}

	if ok { 
		context.JSON(401, gin.H{"error": "token is revoked"})
		context.Abort()
		return
	}

	claims, err := middleware.Verifier.VerifyTokenWithClaims(token)
	if err != nil {
		context.JSON(401, gin.H{"error": "token is invalid"})
		context.Abort()
		return
	}

	log.Info(claims.UserId)

	context.Set("userId", claims.UserId)
	context.Next()
}
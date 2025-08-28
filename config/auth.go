package config

import (
	"time"

	"github.com/Rissochek/messenger-common/utils"
	log "github.com/sirupsen/logrus"
)

type AuthConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

func LoadAuthConfig(logger *log.Logger) *AuthConfig {
	utils.LoadEnvFile()
	secretKey := utils.GetKeyFromEnv("SECRET_KEY")
	tokenDurationStr := utils.GetKeyFromEnv("ACCESS_TOKEN_DURATION")
	
	tokenDurationTime, err := time.ParseDuration(tokenDurationStr) 
	if err != nil {
		logger.Fatalf("failed to parse access duration with error: %v", err)
	}

	return &AuthConfig{SecretKey: secretKey, TokenDuration: tokenDurationTime}
}
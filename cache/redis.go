package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Rissochek/messenger-common/utils"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type BlacklistManager interface {
	AddToBlacklist(token string, expiry int64, ctx context.Context) error
	BlacklistCheck(ctx context.Context, token string) (bool, error)
}

type RedisManager struct {
	RedisClient *redis.Client
	Logger *log.Logger
}

func NewRedisManager(Logger *log.Logger) *RedisManager {
	redis_host := utils.GetKeyFromEnv("REDIS_HOST")
	redis_port := utils.GetKeyFromEnv("REDIS_PORT")
	redis_password := utils.GetKeyFromEnv("REDIS_PASSWORD")
	redis_db_num, err := strconv.Atoi(utils.GetKeyFromEnv("REDIS_DB_NUM"))
	if err != nil {
		Logger.Fatalf("invalid redis_db_num: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", redis_host, redis_port),
		Password: fmt.Sprintf("%v", redis_password),
		DB:       redis_db_num,
	})

	return &RedisManager{RedisClient: client}
}

//checking token without Bearer
func (RedisManager *RedisManager) AddToBlacklist(token string, expiry int64, ctx context.Context) error {
	exparation_time := time.Duration(expiry-time.Now().Unix()) * time.Second
	if err := RedisManager.RedisClient.Set(ctx, token, "revoked", exparation_time); err.Err() != nil {
		RedisManager.Logger.Errorf("failed to add token to blacklist: %v", err.Err())
		return err.Err()
	}

	return nil
}

//bool values is {false: not blacklisted, true: blacklisted}
func (RedisManager *RedisManager) BlacklistCheck(ctx context.Context, token string) (bool, error) {
	result, err := RedisManager.RedisClient.Get(ctx, token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil{
		RedisManager.Logger.Errorf("failed to get token: %v", err)
		return false, fmt.Errorf("failed to check token blacklisted")
	}
	if result == "revoked" {
		RedisManager.Logger.Errorf("token is revoked")
		return true, nil
	}

	return false, nil
}
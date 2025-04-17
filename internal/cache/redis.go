package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

// RedisClient stores the Redis client instance
var RedisClient *redis.Client

// InitializeRedis initializes the Redis client
func InitializeRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // or your Redis server address
		DB:   0,                // default DB
	})

	// Test the connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Could not connect to Redis: %v", err))
	}
}

// SetTokenInCache stores the JWT token in Redis with an expiration time
func SetTokenInCache(token string, expiration time.Duration) error {
	return RedisClient.Set(ctx, token, "valid", expiration).Err()
}

// CheckTokenInCache checks if a JWT token is in the cache (i.e., still valid)
func CheckTokenInCache(token string) (bool, error) {
	result, err := RedisClient.Get(ctx, token).Result()
	if err == redis.Nil {
		// Token not found, it's not in cache
		return false, nil
	} else if err != nil {
		return false, err
	}
	// If "valid" value exists, the token is still valid
	return result == "valid", nil
}

// BlacklistToken adds a JWT token to the blacklist in Redis
func BlacklistToken(token string, expiration time.Duration) error {
	return RedisClient.Set(ctx, "blacklist:"+token, "blacklisted", expiration).Err()
}

// CheckBlacklisted checks if a JWT token is blacklisted
func CheckBlacklisted(token string) (bool, error) {
	_, err := RedisClient.Get(ctx, "blacklist:"+token).Result()
	if err == redis.Nil {
		// Token not found, it's not blacklisted
		return false, nil
	} else if err != nil {
		return false, err
	}
	// If token is found, it's blacklisted
	return true, nil
}

package di

import (
	"database/sql"
	"math-spark/internal/infrastructure/cache"
	"math-spark/internal/infrastructure/email"
	"math-spark/internal/infrastructure/queue"

	"github.com/redis/go-redis/v9"
)

type Container struct {
	// AuthHandler *handlers.AuthHandler
	// UserHandler *handlers.UserHandler
	RedisClient *redis.Client
	CacheService *cache.Cache
}

func NewContainer(db *sql.DB) (*Container, *email.EmailService) {
	// Logger
	// log := logger.NewLogger()


	// Client Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	// Email Queue
	emailQueue := queue.NewQueue(redisClient, "email_queue")

	// Email Service
	emailService := email.NewEmailService(emailQueue)
	
	// Cache Service
	cacheService := cache.NewCache(redisClient)


	return &Container{
		RedisClient: redisClient,
		CacheService: cacheService,
	}, emailService
}

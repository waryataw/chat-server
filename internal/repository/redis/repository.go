package redis

import (
	"github.com/waryataw/chat-server/internal/client/cache"
	"github.com/waryataw/chat-server/internal/service/chat"
)

type repo struct {
	cl cache.RedisClient
}

// NewRepository Конструктор репозитория для кеширования.
func NewRepository(cl cache.RedisClient) chat.AuthCacheRepository {
	return &repo{cl: cl}
}

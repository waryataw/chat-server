package chat

import (
	"github.com/waryataw/chat-server/internal/api/chat"
	"github.com/waryataw/platform_common/pkg/db"
)

type chatService struct {
	authRepository      AuthRepository
	AuthCacheRepository AuthCacheRepository
	repository          Repository
	txManager           db.TxManager
}

// NewService Конструктор Чат сервиса.
func NewService(
	authRepository AuthRepository,
	authCacheRepository AuthCacheRepository,
	chatRepository Repository,
	txManager db.TxManager,
) chat.Service {
	return &chatService{
		authRepository:      authRepository,
		AuthCacheRepository: authCacheRepository,
		repository:          chatRepository,
		txManager:           txManager,
	}
}

package chat

import (
	"github.com/waryataw/auth/pkg/client/db"
	"github.com/waryataw/chat-server/internal/api/chat"
)

type chatService struct {
	authRepository AuthRepository
	repository     Repository
	txManager      db.TxManager
}

// NewService Конструктор Чат сервиса.
func NewService(
	authRepository AuthRepository,
	chatRepository Repository,
	txManager db.TxManager,
) chat.Service {
	return &chatService{
		authRepository: authRepository,
		repository:     chatRepository,
		txManager:      txManager,
	}
}

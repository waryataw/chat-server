package chat

import (
	"github.com/waryataw/chat-server/internal/api/chat"
	"github.com/waryataw/chat-server/internal/client/db"
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

package chat

import (
	"github.com/waryataw/chat-server/internal/repository"
	"github.com/waryataw/chat-server/internal/service"
)

type chatService struct {
	authRepository repository.AuthRepository
	chatRepository repository.ChatRepository
}

// NewService Конструктор Чат сервиса
func NewService(authRepository repository.AuthRepository, chatRepository repository.ChatRepository) service.ChatService {
	return &chatService{authRepository: authRepository, chatRepository: chatRepository}
}

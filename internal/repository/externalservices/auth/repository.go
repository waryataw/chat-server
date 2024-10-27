package auth

import (
	"github.com/waryataw/auth/pkg/authv1"
	"github.com/waryataw/chat-server/internal/service/chat"
)

type repo struct {
	client authv1.AuthServiceClient
}

// NewRepository Конструктор репозитория пользователя.
func NewRepository(client authv1.AuthServiceClient) chat.AuthRepository {
	return &repo{client: client}
}

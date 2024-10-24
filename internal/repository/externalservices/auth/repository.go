package auth

import (
	"github.com/waryataw/auth/pkg/authv1"
	"github.com/waryataw/chat-server/internal/repository"
)

type repo struct {
	client authv1.AuthServiceClient
}

// NewRepository Конструктор репозитория пользователя
func NewRepository(client authv1.AuthServiceClient) repository.AuthRepository {
	return &repo{client: client}
}

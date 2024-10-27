package user

import (
	"github.com/waryataw/chat-server/internal/client/db"
	"github.com/waryataw/chat-server/internal/repository"
)

type repo struct {
	db db.Client
}

// NewRepository Конструктор репозитория пользователя.
func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

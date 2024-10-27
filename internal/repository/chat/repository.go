package user

import (
	"github.com/waryataw/chat-server/internal/client/db"
	"github.com/waryataw/chat-server/internal/service/chat"
)

type repo struct {
	db db.Client
}

// NewRepository Конструктор репозитория пользователя.
func NewRepository(db db.Client) chat.Repository {
	return &repo{db: db}
}

package user

import (
	"github.com/waryataw/chat-server/internal/service/chat"
	"github.com/waryataw/platform_common/pkg/db"
)

type repo struct {
	db db.Client
}

// NewRepository Конструктор репозитория пользователя.
func NewRepository(db db.Client) chat.Repository {
	return &repo{db: db}
}

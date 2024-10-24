package repository

import (
	"context"

	"github.com/waryataw/chat-server/internal/model"
)

// ChatRepository Репозиторий сервиса чата
type ChatRepository interface {
	Create(ctx context.Context, users []*model.User) (int64, error)
	Delete(ctx context.Context, id int64) error
	CreateMessage(ctx context.Context, message *model.Message) error
	Get(ctx context.Context, user *model.User) (*model.Chat, error)
}

// AuthRepository Репозиторий внешнего Auth сервера
type AuthRepository interface {
	GetUser(ctx context.Context, name string) (*model.User, error)
}

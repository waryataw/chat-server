package chat

import (
	"context"

	"github.com/waryataw/chat-server/internal/models"
)

// Repository Репозиторий сервиса чата.
type Repository interface {
	Create(ctx context.Context, users []*models.User) (int64, error)
	Delete(ctx context.Context, id int64) error
	CreateMessage(ctx context.Context, message *models.Message) error
	Get(ctx context.Context, user *models.User) (*models.Chat, error)
}

// AuthRepository Репозиторий внешнего Auth сервера.
type AuthRepository interface {
	GetUser(ctx context.Context, name string) (*models.User, error)
}

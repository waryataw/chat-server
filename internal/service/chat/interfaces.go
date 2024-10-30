package chat

import (
	"context"

	"github.com/waryataw/chat-server/internal/models"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../bin/minimock -i "Repository, AuthRepository, TxManager" -o ./mocks/mocks.go

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

// TxManager Интерфейс моков.
type TxManager interface {
	ReadCommitted(ctx context.Context, f func(ctx context.Context) error) error
}

package chat

import (
	"context"

	"github.com/waryataw/chat-server/internal/models"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../bin/minimock -i "Repository,AuthRepository,AuthCacheRepository,github.com/waryataw/platform_common/pkg/db.TxManager" -o ./mocks/ -s "_minimock.go"

// Repository Интерфейс репозитория сервиса чата.
type Repository interface {
	Create(ctx context.Context, users []*models.User) (int64, error)
	Delete(ctx context.Context, id int64) error
	CreateMessage(ctx context.Context, message *models.Message) error
	Get(ctx context.Context, user *models.User) (*models.Chat, error)
}

// AuthRepository Интерфейс репозитория внешнего Auth сервера.
type AuthRepository interface {
	GetUser(ctx context.Context, name string) (*models.User, error)
}

// AuthCacheRepository Интерфейс репозитория для кеша.
type AuthCacheRepository interface {
	GetUser(ctx context.Context, name string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

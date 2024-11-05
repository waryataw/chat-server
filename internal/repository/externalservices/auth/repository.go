package auth

import (
	"github.com/waryataw/chat-server/internal/service/chat"
	"github.com/waryataw/platform_common/pkg/authclient"
)

type repo struct {
	client *authclient.AuthClient
}

// NewRepository Конструктор репозитория пользователя.
func NewRepository(client *authclient.AuthClient) chat.AuthRepository {
	return &repo{client: client}
}

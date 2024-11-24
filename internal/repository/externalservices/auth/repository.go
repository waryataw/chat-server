package auth

import (
	"github.com/waryataw/chat-server/internal/service/chat"
	"github.com/waryataw/platform_common/pkg/userclient"
)

type repo struct {
	client *userclient.UserClient
}

// NewRepository Конструктор репозитория пользователя.
func NewRepository(client *userclient.UserClient) chat.AuthRepository {
	return &repo{client: client}
}

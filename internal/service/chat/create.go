package chat

import (
	"context"
	"fmt"

	"github.com/waryataw/chat-server/internal/model"
)

func (s *chatService) Create(ctx context.Context, usernames []string) (int64, error) {
	users := make([]*model.User, len(usernames))
	for index, username := range usernames {
		user, err := s.authRepository.GetUser(ctx, username)
		if err != nil {
			return 0, fmt.Errorf("failed get user from auth service: %w", err)
		}

		users[index] = user
	}

	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.Create(ctx, users)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed create chat: %w", err)
	}

	return id, nil
}
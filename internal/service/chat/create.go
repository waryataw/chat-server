package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/waryataw/chat-server/internal/models"
)

func (s chatService) Create(ctx context.Context, usernames []string) (int64, error) {
	users := make([]*models.User, len(usernames))
	for index, username := range usernames {
		user, err := s.getUser(ctx, username)
		if err != nil {
			return 0, fmt.Errorf("failet get user: %w", err)
		}

		users[index] = user
	}

	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		id, err = s.repository.Create(ctx, users)

		return err
	})

	if err != nil {
		return 0, fmt.Errorf("failed create chat: %w", err)
	}

	return id, nil
}

func (s chatService) getUser(ctx context.Context, username string) (*models.User, error) {
	user, err := s.AuthCacheRepository.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			user, err = s.authRepository.GetUser(ctx, username)
			if err != nil {
				return nil, fmt.Errorf("failed to get user from auth service: %w", err)
			}

			if err = s.AuthCacheRepository.CreateUser(ctx, user); err != nil {
				return nil, fmt.Errorf("failed to create user in cache: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to get user from auth cache: %w", err)
	}

	return user, nil
}

package redis

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/waryataw/chat-server/internal/models"
)

func (r repo) GetUser(ctx context.Context, name string) (*models.User, error) {
	values, err := r.cl.HGetAll(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if len(values) == 0 {
		return nil, models.ErrUserNotFound
	}

	var user models.User
	if err = redigo.ScanStruct(values, &user); err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

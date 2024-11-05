package redis

import (
	"context"
	"fmt"

	"github.com/waryataw/chat-server/internal/models"
)

func (r repo) CreateUser(ctx context.Context, user *models.User) error {
	if err := r.cl.HashSet(ctx, user.Name, user); err != nil {
		return fmt.Errorf("failed user hash set: %w", err)
	}

	return nil
}

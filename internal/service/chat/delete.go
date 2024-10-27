package chat

import (
	"context"
	"fmt"
)

func (s chatService) Delete(ctx context.Context, id int64) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete chat %d: %w", id, err)
	}

	return nil
}

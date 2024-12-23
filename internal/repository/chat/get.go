package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/waryataw/chat-server/internal/models"
	"github.com/waryataw/platform_common/pkg/db"
)

// Get Метод получения чата для пользователя.
func (r repo) Get(ctx context.Context, user *models.User) (*models.Chat, error) {
	// Пока выберу первый попавшийся, потом будет совсем иначе все.
	builderSelect := sq.
		Select(
			"chat.id",
			"chat.created_at",
		).
		From("chat").
		Join("chat_user ON chat_user.chat_id = chat.id").
		Where(sq.Eq{"chat_user.user_id": user.ID}).
		Limit(1)

	sql, args, err := builderSelect.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build chat selection query: %w", err)
	}

	query := db.Query{
		Name:     "chat_repository.Get",
		QueryRaw: sql,
	}

	var chat models.Chat
	if err := r.db.DB().QueryRowContext(ctx, query, args...).Scan(&chat.ID, &chat.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no chat found for user: %d", user.ID)
		}
		return nil, fmt.Errorf("failed to select chat: %w", err)
	}

	return &chat, nil
}

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

// CreateMessage Метод создания сообщения.
func (r repo) CreateMessage(ctx context.Context, message *models.Message) error {
	// Пока выберу первый попавшийся, потом будет совсем иначе все.
	builderSelect := sq.Select("chat_id").
		From("chat_user").
		Where(sq.Eq{"user_id": message.User.ID}).
		Limit(1)

	sql, args, err := builderSelect.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build chat selection query: %w", err)
	}

	query := db.Query{
		Name:     "chat_repository.CreateMessage.GetChat",
		QueryRaw: sql,
	}

	var chatID int64
	if err := r.db.DB().QueryRowContext(ctx, query, args...).Scan(&chatID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("no chat found for user: %d", message.User.ID)
		}
		return fmt.Errorf("failed to select chat: %w", err)
	}

	builderInsert := sq.Insert("message").
		Columns(
			"chat_id",
			"user_id",
			"text",
		).
		Values(
			chatID,
			message.User.ID,
			message.Text,
		)

	sql, args, err = builderInsert.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build message insertation query: %w", err)
	}

	queryInsert := db.Query{
		Name:     "chat_repository.CreateMessage",
		QueryRaw: sql,
	}

	if _, err := r.db.DB().ExecContext(ctx, queryInsert, args...); err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	return nil
}

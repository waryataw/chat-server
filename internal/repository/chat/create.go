package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/waryataw/chat-server/internal/client/db"
	"github.com/waryataw/chat-server/internal/model"
)

func (r *repo) Create(ctx context.Context, users []*model.User) (int64, error) {
	builder := sq.Insert("chat").
		Columns("created_at").
		Values(sq.Expr("NOW()")).
		Suffix("RETURNING id")

	sql, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query to insert chat: %w", err)
	}

	query := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: sql,
	}

	var chatID int64
	if err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&chatID); err != nil {
		return 0, fmt.Errorf("failed to insert chat: %w", err)
	}

	for _, user := range users {
		builderInsert := sq.Insert("chat_user").
			Columns(
				"chat_id",
				"user_id",
			).
			Values(chatID, user.ID)

		sql, args, err = builderInsert.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return 0, fmt.Errorf("failed to build query to insert chat_user: %w", err)
		}

		queryInsertChatUser := db.Query{
			Name:     "chat_repository.CreateChatUser",
			QueryRaw: sql,
		}

		_, err = r.db.DB().ExecContext(ctx, queryInsertChatUser, args...)
		if err != nil {
			return 0, fmt.Errorf("failed to insert chat_user: %w", err)
		}
	}

	return chatID, nil
}
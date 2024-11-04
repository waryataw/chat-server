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

	// Вот тут вопрос, ошибку отдавать если не найдено или nil, nil. Мы же понимаем что отсутствие пользователя в кеше
	// это не ошибка, а он просто еще туда не записался. Пока остановился на таком варианте.
	if len(values) == 0 {
		return nil, nil
	}

	var user models.User
	if err = redigo.ScanStruct(values, &user); err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

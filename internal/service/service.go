package service

import (
	"context"
)

// ChatService Интерфейс сервиса Чата
type ChatService interface {
	Create(ctx context.Context, usernames []string) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, from string, text string) error
}

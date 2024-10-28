package chat

import "context"

// Service Интерфейс сервиса Чата.
type Service interface {
	Create(ctx context.Context, usernames []string) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, from string, text string) error
}

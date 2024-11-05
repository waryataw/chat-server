package chat

import "context"

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../bin/minimock -i github.com/waryataw/chat-server/internal/api/chat.* -o "./mocks/mocks.go"

// Service Интерфейс сервиса Чата.
type Service interface {
	Create(ctx context.Context, usernames []string) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, from string, text string) error
}

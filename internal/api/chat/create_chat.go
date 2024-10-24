package chat

import (
	"context"
	"fmt"

	"github.com/waryataw/chat-server/pkg/chatserverv1"
)

// CreateChat Метод создания Чата
func (i *Implementation) CreateChat(ctx context.Context, req *chatserverv1.CreateChatRequest) (*chatserverv1.CreateChatResponse, error) {
	id, err := i.chatService.Create(ctx, req.Usernames)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	return &chatserverv1.CreateChatResponse{
		Id: id,
	}, nil
}

package chat

import (
	"context"
	"fmt"

	"github.com/waryataw/chat-server/internal/model"
)

func (s *chatService) SendMessage(ctx context.Context, from string, text string) error {
	user, err := s.authRepository.GetUser(ctx, from)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	chat, err := s.chatRepository.Get(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to get chat: %w", err)
	}

	message := model.Message{
		Chat: chat,
		User: user,
		Text: text,
	}

	if err := s.chatRepository.CreateMessage(ctx, &message); err != nil {
		return fmt.Errorf("failed to create chat message: %w", err)
	}

	return nil
}

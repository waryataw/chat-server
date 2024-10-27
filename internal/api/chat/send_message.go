package chat

import (
	"context"
	"fmt"

	"github.com/waryataw/chat-server/pkg/chatserverv1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SendMessage Метод отправки сообщения в чат.
func (c Controller) SendMessage(ctx context.Context, req *chatserverv1.SendMessageRequest) (*emptypb.Empty, error) {
	if err := c.chatService.SendMessage(ctx, req.GetFrom(), req.GetText()); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return &emptypb.Empty{}, nil
}

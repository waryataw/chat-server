package chat

import (
	"context"
	"fmt"

	"github.com/waryataw/chat-server/pkg/chatserverv1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteChat Метод удаления чата.
func (c *Controller) DeleteChat(ctx context.Context, req *chatserverv1.DeleteChatRequest) (*emptypb.Empty, error) {
	if err := c.chatService.Delete(ctx, req.GetId()); err != nil {
		return nil, fmt.Errorf("failed to delete chat: %w", err)
	}

	return &emptypb.Empty{}, nil
}

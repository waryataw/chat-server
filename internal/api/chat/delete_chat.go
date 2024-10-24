package chat

import (
	"context"
	"fmt"

	"github.com/waryataw/chat-server/pkg/chatserverv1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteChat Метод удаления чата
func (i *Implementation) DeleteChat(ctx context.Context, req *chatserverv1.DeleteChatRequest) (*emptypb.Empty, error) {
	if err := i.chatService.Delete(ctx, req.GetId()); err != nil {
		return nil, fmt.Errorf("failed to delete chat: %w", err)
	}

	return &emptypb.Empty{}, nil
}

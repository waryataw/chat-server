package chat

import (
	"github.com/waryataw/chat-server/internal/service"
	"github.com/waryataw/chat-server/pkg/chatserverv1"
)

// Controller Реализация Чат сервиса.
type Controller struct {
	chatserverv1.UnimplementedChatServerServiceServer
	chatService service.ChatService
}

// NewController Конструктор реализации Чат сервиса.
func NewController(chatService service.ChatService) *Controller {
	return &Controller{chatService: chatService}
}

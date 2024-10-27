package chat

import (
	"github.com/waryataw/chat-server/internal/service"
	"github.com/waryataw/chat-server/pkg/chatserverv1"
)

// Implementation Реализация Чат сервиса.
type Implementation struct {
	chatserverv1.UnimplementedChatServerServiceServer
	chatService service.ChatService
}

// NewImplementation Конструктор реализации Чат сервиса.
func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{chatService: chatService}
}

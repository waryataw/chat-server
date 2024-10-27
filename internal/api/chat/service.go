package chat

import (
	"github.com/waryataw/chat-server/pkg/chatserverv1"
)

// Controller Реализация Чат сервиса.
type Controller struct {
	chatserverv1.UnimplementedChatServerServiceServer
	service Service
}

// NewController Конструктор реализации Чат сервиса.
func NewController(service Service) *Controller {
	return &Controller{service: service}
}

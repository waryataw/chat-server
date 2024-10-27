package models

import "time"

// Chat Чат.
type Chat struct {
	ID        int64
	CreatedAt *time.Time
	Users     []*User
}

// Message Сообщение.
type Message struct {
	ID        int64
	Chat      *Chat
	User      *User
	Text      string
	CreatedAt *time.Time
}

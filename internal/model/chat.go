package model

import "time"

// User Пользователь
type User struct {
	ID int64
}

// Chat Чат
type Chat struct {
	ID        int64
	CreatedAt *time.Time
}

// ChatUser Сущность связи Чата к пользователям (А может она и не нужна, вопрос, наверное chat должен содержать users)
type ChatUser struct {
	ID   int64
	Chat *Chat
	User *User
}

// Message Сообщение
type Message struct {
	ID        int64
	Chat      *Chat
	User      *User
	Text      string
	CreatedAt *time.Time
}

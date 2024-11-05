package models

import "github.com/pkg/errors"

// ErrUserNotFound Ошибка: Пользователь не найден.
var ErrUserNotFound = errors.New("user not found")

// User Пользователь.
type User struct {
	ID   int64
	Name string
}

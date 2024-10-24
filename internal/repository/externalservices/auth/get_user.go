package auth

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/pkg/authv1"
	"github.com/waryataw/chat-server/internal/model"
)

func (repo *repo) GetUser(ctx context.Context, name string) (*model.User, error) {
	user, err := repo.client.GetUser(ctx, &authv1.GetUserRequest{
		Query: &authv1.GetUserRequest_Name{
			Name: name,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed get user from auth service: %w", err)
	}

	return &model.User{ID: user.GetId()}, nil
}

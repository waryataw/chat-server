package auth

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/pkg/userv1"
	"github.com/waryataw/chat-server/internal/models"
)

func (r repo) GetUser(ctx context.Context, name string) (*models.User, error) {
	user, err := r.client.UserServiceClient.GetUser(ctx, &userv1.GetUserRequest{
		Query: &userv1.GetUserRequest_Name{
			Name: name,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed get user from auth service: %w", err)
	}

	return &models.User{ID: user.GetId(), Name: name}, nil
}

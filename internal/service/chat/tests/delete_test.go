package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/waryataw/chat-server/internal/service/chat"
	"github.com/waryataw/chat-server/internal/service/chat/mocks"
	"github.com/waryataw/platform_common/pkg/db"
)

func TestDelete(t *testing.T) {
	type repositoryMockBehavior func(mc *minimock.Controller) chat.Repository
	type authRepositoryMockBehavior func(mc *minimock.Controller) chat.AuthRepository
	type txManagerMockBehavior func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repoError = fmt.Errorf("failed to delete chat %d", id)
	)

	tests := []struct {
		name                       string
		args                       args
		want                       int64
		err                        error
		repositoryMockBehavior     repositoryMockBehavior
		authRepositoryMockBehavior authRepositoryMockBehavior
		txManagerMockBehavior      txManagerMockBehavior
	}{
		{
			"success case",
			args{
				ctx: ctx,
				id:  id,
			},
			0,
			nil,
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)

				return mock
			},
		},
		{
			"repo error case",
			args{
				ctx: ctx,
				id:  id,
			},
			0,
			fmt.Errorf("failed to delete chat %d: %w", id, repoError),
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(repoError)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repositoryMock := tt.repositoryMockBehavior(mc)
			authRepositoryMock := tt.authRepositoryMockBehavior(mc)
			txManagerMock := tt.txManagerMockBehavior(mc)
			service := chat.NewService(authRepositoryMock, repositoryMock, txManagerMock)

			err := service.Delete(tt.args.ctx, tt.args.id)

			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}
}

package tests

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/waryataw/chat-server/internal/models"
	"github.com/waryataw/chat-server/internal/service/chat"
	"github.com/waryataw/chat-server/internal/service/chat/mocks"
	"github.com/waryataw/platform_common/pkg/db"
	"testing"
)

func TestCreate(t *testing.T) {
	type repositoryMockBehavior func(mc *minimock.Controller) chat.Repository
	type authRepositoryMockBehavior func(mc *minimock.Controller) chat.AuthRepository
	type txManagerMockBehavior func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx       context.Context
		usernames []string
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		//serviceErr = fmt.Errorf("failed to create chat")

		id        = gofakeit.Int64()
		user      = &models.User{ID: gofakeit.Int64()}
		users     = []*models.User{user}
		usernames = []string{gofakeit.Username()}
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
				ctx:       ctx,
				usernames: usernames,
			},
			id,
			nil,
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, users).Return(id, nil)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, usernames[0]).Return(user, nil)

				return mock
			},
			func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Expect(ctx, func(_ context.Context) error {
					return nil
				}).Return(nil)

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
			api := chat.NewService(authRepositoryMock, repositoryMock, txManagerMock)

			response, err := api.Create(tt.args.ctx, tt.args.usernames)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}

}

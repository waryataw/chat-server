package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/waryataw/chat-server/internal/models"
	"github.com/waryataw/chat-server/internal/service/chat"
	"github.com/waryataw/chat-server/internal/service/chat/mocks"
	"github.com/waryataw/platform_common/pkg/db"
)

func TestCreate(t *testing.T) {
	type repositoryMockBehavior func(mc *minimock.Controller) chat.Repository
	type authRepositoryMockBehavior func(mc *minimock.Controller) chat.AuthRepository
	type authCacheRepositoryMockBehavior func(mc *minimock.Controller) chat.AuthCacheRepository
	type txManagerMockBehavior func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx       context.Context
		usernames []string
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		user      = &models.User{ID: gofakeit.Int64()}
		usernames = []string{gofakeit.Username()}

		userFromCacheErr       = fmt.Errorf("failed get user from cache")
		userFromAuthServiceErr = fmt.Errorf("failed get user from auth service")
		userToCacheErr         = fmt.Errorf("failed to create user in cache")
		txManagerErr           = fmt.Errorf("tx commit failed")
	)

	tests := []struct {
		name                            string
		args                            args
		want                            int64
		err                             error
		repositoryMockBehavior          repositoryMockBehavior
		authRepositoryMockBehavior      authRepositoryMockBehavior
		authCacheRepositoryMockBehavior authCacheRepositoryMockBehavior
		txManagerMockBehavior           txManagerMockBehavior
	}{
		{
			"success case",
			args{
				ctx:       ctx,
				usernames: usernames,
			},
			0,
			nil,
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthCacheRepository {
				mock := mocks.NewAuthCacheRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, usernames[0]).Return(user, nil)

				return mock
			},
			func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.ExpectCtxParam1(ctx).Return(nil)

				return mock
			},
		},
		{
			"user from cache error case",
			args{
				ctx:       ctx,
				usernames: usernames,
			},
			0,
			fmt.Errorf("failet get user: failed to get user from auth cache: %w", userFromCacheErr),
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthCacheRepository {
				mock := mocks.NewAuthCacheRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, usernames[0]).Return(nil, userFromCacheErr)

				return mock
			},
			func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)

				return mock
			},
		},
		{
			"user from auth error case",
			args{
				ctx:       ctx,
				usernames: usernames,
			},
			0,
			fmt.Errorf("failet get user: failed to get user from auth service: %w", userFromAuthServiceErr),
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, usernames[0]).Return(nil, userFromAuthServiceErr)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthCacheRepository {
				mock := mocks.NewAuthCacheRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, usernames[0]).Return(nil, models.ErrUserNotFound)

				return mock
			},
			func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)

				return mock
			},
		},
		{
			"user to cache error case",
			args{
				ctx:       ctx,
				usernames: usernames,
			},
			0,
			fmt.Errorf("failet get user: failed to create user in cache: %w", userToCacheErr),
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, usernames[0]).Return(user, nil)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthCacheRepository {
				mock := mocks.NewAuthCacheRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, usernames[0]).Return(nil, models.ErrUserNotFound)
				mock.CreateUserMock.Expect(ctx, user).Return(userToCacheErr)

				return mock
			},
			func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)

				return mock
			},
		},
		{
			"tx manager error",
			args{
				ctx:       ctx,
				usernames: usernames,
			},
			0,
			fmt.Errorf("failed create chat: %w", txManagerErr),
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthCacheRepository {
				mock := mocks.NewAuthCacheRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, usernames[0]).Return(user, nil)

				return mock
			},
			func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.ExpectCtxParam1(ctx).Return(txManagerErr)

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
			authCacheRepositoryMock := tt.authCacheRepositoryMockBehavior(mc)
			txManagerMock := tt.txManagerMockBehavior(mc)
			service := chat.NewService(
				authRepositoryMock,
				authCacheRepositoryMock,
				repositoryMock,
				txManagerMock,
			)

			response, err := service.Create(tt.args.ctx, tt.args.usernames)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}
}

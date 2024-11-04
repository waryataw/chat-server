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

func TestSendMessage(t *testing.T) {
	type repositoryMockBehavior func(mc *minimock.Controller) chat.Repository
	type authRepositoryMockBehavior func(mc *minimock.Controller) chat.AuthRepository
	type authCacheRepositoryMockBehavior func(mc *minimock.Controller) chat.AuthCacheRepository
	type txManagerMockBehavior func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx       context.Context
		username  string
		message   *models.Message
		modelChat *models.Chat
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		userID   = gofakeit.Int64()
		username = gofakeit.Username()
		user     = &models.User{ID: userID}

		modelChat = &models.Chat{ID: gofakeit.Int64(), Users: []*models.User{user}}

		message = &models.Message{
			ID:   gofakeit.Int64(),
			Chat: modelChat,
			User: user,
			Text: gofakeit.LoremIpsumSentence(5),
		}

		authRepoErr              = fmt.Errorf("failed get user from auth service")
		repoGetUserErr           = fmt.Errorf("failed to get chat")
		repoCreateChatMessageErr = fmt.Errorf("failed to create chat message")
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
				username:  username,
				message:   message,
				modelChat: modelChat,
			},
			0,
			nil,
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.GetMock.Expect(ctx, user).Return(modelChat, nil)
				mock.CreateMessageMock.Set(func(_ context.Context, message *models.Message) (err error) {
					require.Equal(t, modelChat.ID, message.Chat.ID)
					return nil
				})

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, username).Return(user, nil)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthCacheRepository {
				mock := mocks.NewAuthCacheRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)

				return mock
			},
		},
		{
			"auth repo error case",
			args{
				ctx:       ctx,
				username:  username,
				message:   message,
				modelChat: modelChat,
			},
			0,
			fmt.Errorf("failed to get user: %w", authRepoErr),
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, username).Return(nil, authRepoErr)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthCacheRepository {
				mock := mocks.NewAuthCacheRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)

				return mock
			},
		},
		{
			"repo get chat error case",
			args{
				ctx:       ctx,
				username:  username,
				message:   message,
				modelChat: modelChat,
			},
			0,
			fmt.Errorf("failed to get chat: %w", repoGetUserErr),
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.GetMock.Expect(ctx, user).Return(modelChat, repoGetUserErr)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, username).Return(user, nil)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthCacheRepository {
				mock := mocks.NewAuthCacheRepositoryMock(mc)

				return mock
			},
			func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)

				return mock
			},
		},
		{
			"repo create chat message error case",
			args{
				ctx:       ctx,
				username:  username,
				message:   message,
				modelChat: modelChat,
			},
			0,
			fmt.Errorf("failed to create chat message: %w", repoCreateChatMessageErr),
			func(mc *minimock.Controller) chat.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.GetMock.Expect(ctx, user).Return(modelChat, nil)
				mock.CreateMessageMock.Set(func(_ context.Context, message *models.Message) (err error) {
					require.Equal(t, modelChat.ID, message.Chat.ID)
					return repoCreateChatMessageErr
				})

				return mock
			},
			func(mc *minimock.Controller) chat.AuthRepository {
				mock := mocks.NewAuthRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, username).Return(user, nil)

				return mock
			},
			func(mc *minimock.Controller) chat.AuthCacheRepository {
				mock := mocks.NewAuthCacheRepositoryMock(mc)

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
			authCacheRepositoryMock := tt.authCacheRepositoryMockBehavior(mc)
			txManagerMock := tt.txManagerMockBehavior(mc)
			service := chat.NewService(
				authRepositoryMock,
				authCacheRepositoryMock,
				repositoryMock,
				txManagerMock,
			)

			err := service.SendMessage(tt.args.ctx, tt.args.username, tt.args.message.Text)

			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}
}

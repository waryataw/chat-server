package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/waryataw/chat-server/internal/api/chat"
	"github.com/waryataw/chat-server/internal/api/chat/mocks"
	"github.com/waryataw/chat-server/pkg/chatserverv1"
	"testing"
)

func TestCreateChat(t *testing.T) {
	type mockBehavior func(mc *minimock.Controller) chat.Service

	type args struct {
		ctx context.Context
		req *chatserverv1.CreateChatRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("failed to create chat")

		id = gofakeit.Int64()

		req = &chatserverv1.CreateChatRequest{
			Usernames: []string{
				gofakeit.Username(),
				gofakeit.Username(),
			},
		}
		res = &chatserverv1.CreateChatResponse{Id: id}
	)

	tests := []struct {
		name         string
		args         args
		want         *chatserverv1.CreateChatResponse
		err          error
		mockBehavior mockBehavior
	}{
		{
			"success case",
			args{
				ctx: ctx,
				req: req,
			},
			res,
			nil,
			func(mc *minimock.Controller) chat.Service {
				mock := mocks.NewServiceMock(mc)
				mock.CreateMock.Expect(ctx, req.Usernames).Return(id, nil)

				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  fmt.Errorf("failed to create chat: %w", serviceErr),
			mockBehavior: func(mc *minimock.Controller) chat.Service {
				mock := mocks.NewServiceMock(mc)
				mock.CreateMock.Expect(ctx, req.Usernames).Return(0, serviceErr)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := tt.mockBehavior(mc)
			api := chat.NewController(mock)

			response, err := api.CreateChat(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}
}

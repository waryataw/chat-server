package interceptors

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/pkg/accessv1"
	"github.com/waryataw/platform_common/pkg/accessclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor Auth интерцептор.
type AuthInterceptor struct {
	accessClient *accessclient.AccessClient
}

// NewAuthInterceptor Конструктор Auth интерцептора.
func NewAuthInterceptor(client *accessclient.AccessClient) *AuthInterceptor {
	return &AuthInterceptor{
		accessClient: client,
	}
}

// UnaryInterceptor Прокидывает токен. Проверяет доступ текущего пользователя.
func (a *AuthInterceptor) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is required")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization header is required")
	}

	outgoingCtx := metadata.AppendToOutgoingContext(ctx, "authorization", authHeader[0])

	_, err := a.accessClient.AccessClient.Check(outgoingCtx, &accessv1.CheckRequest{
		EndpointAddress: info.FullMethod,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check access: %w", err)
	}

	return handler(ctx, req)
}

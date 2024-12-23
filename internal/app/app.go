package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/waryataw/chat-server/internal/config"
	interceptors "github.com/waryataw/chat-server/internal/interceptor"
	"github.com/waryataw/chat-server/pkg/chatserverv1"
	"github.com/waryataw/platform_common/pkg/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App Приложение.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp Конструктор приложения.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init deps: %w", err)
	}

	return a, nil
}

// Run Запуск GRPC сервера.
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(
			interceptors.NewAuthInterceptor(a.serviceProvider.AccessClient(ctx)).UnaryInterceptor,
		),
	)

	reflection.Register(a.grpcServer)

	chatserverv1.RegisterChatServerServiceServer(a.grpcServer, a.serviceProvider.ChatController(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	listener, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	err = a.grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

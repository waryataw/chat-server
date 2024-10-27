package app

import (
	"context"
	"log"

	"github.com/waryataw/auth/pkg/authv1"
	"github.com/waryataw/chat-server/internal/client/db/transaction"
	"github.com/waryataw/chat-server/internal/repository"
	chatRepository "github.com/waryataw/chat-server/internal/repository/chat"
	authRepository "github.com/waryataw/chat-server/internal/repository/externalservices/auth"
	"github.com/waryataw/chat-server/internal/service"
	chatService "github.com/waryataw/chat-server/internal/service/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/waryataw/chat-server/internal/api/chat"

	"github.com/waryataw/chat-server/internal/client/db"
	"github.com/waryataw/chat-server/internal/client/db/pg"
	"github.com/waryataw/chat-server/internal/closer"
	"github.com/waryataw/chat-server/internal/config"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	chatRepository repository.ChatRepository

	authClient     authv1.AuthServiceClient
	authRepository repository.AuthRepository

	chatService service.ChatService

	chatController *chat.Controller
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) AuthClient(_ context.Context) authv1.AuthServiceClient {
	if s.authClient == nil {
		grpcClientConfig, err := config.NewGRPCClientConfig()
		if err != nil {
			log.Fatalf("failed to get grpc client config: %v", err)
		}

		conn, err := grpc.NewClient(
			grpcClientConfig.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("failed to connect to Auth server: %v", err)
		}

		s.authClient = authv1.NewAuthServiceClient(conn)
		closer.Add(conn.Close)
	}

	return s.authClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.AuthClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}
	return s.chatRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.AuthRepository(ctx),
			s.ChatRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) ChatController(ctx context.Context) *chat.Controller {
	if s.chatController == nil {
		s.chatController = chat.NewController(s.ChatService(ctx))
	}

	return s.chatController
}

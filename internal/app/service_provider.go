package app

import (
	"context"
	"log"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/waryataw/chat-server/internal/api/chat"
	"github.com/waryataw/chat-server/internal/client/cache"
	"github.com/waryataw/chat-server/internal/client/cache/redis"
	"github.com/waryataw/chat-server/internal/config"
	"github.com/waryataw/chat-server/internal/config/env"
	chatRepository "github.com/waryataw/chat-server/internal/repository/chat"
	authRepository "github.com/waryataw/chat-server/internal/repository/externalservices/auth"
	redisRepo "github.com/waryataw/chat-server/internal/repository/redis"
	chatService "github.com/waryataw/chat-server/internal/service/chat"
	"github.com/waryataw/platform_common/pkg/closer"
	"github.com/waryataw/platform_common/pkg/db"
	"github.com/waryataw/platform_common/pkg/db/pg"
	"github.com/waryataw/platform_common/pkg/db/transaction"
	"github.com/waryataw/platform_common/pkg/userclient"
)

type serviceProvider struct {
	pgConfig    env.PGConfig
	grpcConfig  env.GRPCConfig
	redisConfig config.RedisConfig

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	chatRepository chatService.Repository

	userClient          *userclient.UserClient
	authRepository      chatService.AuthRepository
	authCacheRepository chatService.AuthCacheRepository

	chatService chat.Service

	chatController *chat.Controller
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() env.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
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

func (s *serviceProvider) AuthClient(_ context.Context) *userclient.UserClient {
	if s.userClient == nil {
		grpcClientConfig, err := env.NewGRPCClientConfig()
		if err != nil {
			log.Fatalf("failed to get grpc client config: %v", err)
		}

		authClient, err := userclient.New(grpcClientConfig.Address())
		if err != nil {
			log.Fatalf("failed to create auth client: %v", err)
		}

		s.userClient = authClient
		closer.Add(s.userClient.Conn.Close)
	}

	return s.userClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AuthRepository(ctx context.Context) chatService.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.AuthClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

func (s *serviceProvider) AuthCacheRepository(_ context.Context) chatService.AuthCacheRepository {
	if s.authCacheRepository == nil {
		s.authCacheRepository = redisRepo.NewRepository(s.RedisClient())
	}

	return s.authCacheRepository
}

func (s *serviceProvider) ChatRepository(ctx context.Context) chatService.Repository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) chat.Service {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.AuthRepository(ctx),
			s.AuthCacheRepository(ctx),
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

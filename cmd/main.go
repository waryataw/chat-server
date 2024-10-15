package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"log"
	"net"

	"github.com/waryataw/chat-server/internal/config"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	sq "github.com/Masterminds/squirrel"

	authv1 "github.com/waryataw/auth/pkg/authv1"
	chatserverv1 "github.com/waryataw/chat-server/pkg/chatserverv1"
)

const errNoRows = "no rows in result set"

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	chatserverv1.UnimplementedChatServerServiceServer
	pool       *pgxpool.Pool
	authClient authv1.AuthServiceClient
}

// CreateChat Создание чата
func (s *server) CreateChat(ctx context.Context, req *chatserverv1.CreateChatRequest) (*chatserverv1.CreateChatResponse, error) {
	usernames := req.GetUsernames()
	userIDs := make([]int64, len(usernames))

	for index, username := range usernames {
		user, err := s.authClient.GetUser(ctx, &authv1.GetUserRequest{
			Query: &authv1.GetUserRequest_Name{
				Name: username,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed get user from auth service: %w", err)
		}

		userIDs[index] = user.GetId()
	}

	query := sq.Insert("chat").
		Columns("created_at").
		Values(sq.Expr("NOW()")).
		Suffix("RETURNING id")

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query to insert chat: %w", err)
	}

	var chatID int64
	if err := s.pool.QueryRow(ctx, sql, args...).Scan(&chatID); err != nil {
		return nil, fmt.Errorf("failed to insert chat: %w", err)
	}

	log.Printf("inserted chat with id: %d", chatID)

	for _, id := range userIDs {
		query := sq.Insert("chat_user").
			Columns(
				"chat_id",
				"user_id",
			).
			Values(chatID, id)

		sql, args, err = query.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return nil, fmt.Errorf("failed to build query to insert chat_user: %w", err)
		}

		if _, err := s.pool.Exec(ctx, sql, args...); err != nil {
			return nil, fmt.Errorf("failed to insert chat_user: %w", err)
		}
	}

	return &chatserverv1.CreateChatResponse{Id: chatID}, nil
}

// DeleteChat Удаление чата
func (s *server) DeleteChat(ctx context.Context, req *chatserverv1.DeleteChatRequest) (*emptypb.Empty, error) {
	queryDelete := sq.Delete("chat").Where(sq.Eq{"id": req.GetId()})

	sql, args, err := queryDelete.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query to delete chat: %w", err)
	}

	tag, err := s.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("failed to delete chat: %d not found", req.GetId())
	}

	return &emptypb.Empty{}, nil
}

// SendMessage Отправка сообщения
func (s *server) SendMessage(ctx context.Context, req *chatserverv1.SendMessageRequest) (*emptypb.Empty, error) {
	user, err := s.authClient.GetUser(ctx, &authv1.GetUserRequest{
		Query: &authv1.GetUserRequest_Name{Name: req.GetFrom()},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user from authv1: %w", err)
	}

	// Пока выберу первый попавшийся, потом будет совсем иначе все
	querySelect := sq.Select("chat_id").
		From("chat_user").
		Where(sq.Eq{"user_id": user.GetId()}).
		Limit(1)

	sql, args, err := querySelect.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var chatID int64
	if err := s.pool.QueryRow(ctx, sql, args...).Scan(&chatID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no chat found for user: %d", user.GetId())
		}
		return nil, fmt.Errorf("failed to select chat: %w", err)
	}

	queryInsert := sq.Insert("message").
		Columns(
			"chat_id",
			"user_id",
			"text",
		).
		Values(chatID, user.GetId(), req.GetText())

	sql, args, err = queryInsert.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query to insert message: %w", err)
	}

	if _, err := s.pool.Exec(ctx, sql, args...); err != nil {
		return nil, fmt.Errorf("failed to insert message: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

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
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("failed to close grpc connection: %v", err)
		}
	}(conn)

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)

	client := authv1.NewAuthServiceClient(conn)

	chatserverv1.RegisterChatServerServiceServer(s, &server{pool: pool, authClient: client})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

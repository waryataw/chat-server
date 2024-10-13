package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/waryataw/chat-server/internal/config"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	sq "github.com/Masterminds/squirrel"

	descUser "github.com/waryataw/auth/pkg/auth_v1"
	desc "github.com/waryataw/chat-server/pkg/chat_server_v1"
)

const errNoRows = "no rows in result set"

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatServerServiceServer
	pool       *pgxpool.Pool
	grpcClient descUser.AuthServiceClient
}

func (s *server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {

	var userIDs []int64

	for _, username := range req.GetUsernames() {

		user, err := s.grpcClient.GetUser(ctx, &descUser.GetUserRequest{
			Query: &descUser.GetUserRequest_Name{
				Name: username,
			},
		})
		if err != nil {
			return nil, err
		}

		userIDs = append(userIDs, user.GetId())

	}

	query := sq.Insert("chat").
		PlaceholderFormat(sq.Dollar).
		Columns("created_at").
		Values(sq.Expr("NOW()")).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query to insert chat: %v", err)
	}

	var chatID int64
	err = s.pool.QueryRow(ctx, sql, args...).Scan(&chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert chat: %v", err)
	}

	log.Printf("inserted chat with id: %d", chatID)

	for _, id := range userIDs {

		query := sq.Insert("chat_user").
			PlaceholderFormat(sq.Dollar).
			Columns("chat_id", "user_id").
			Values(chatID, id)

		sql, args, err = query.ToSql()
		if err != nil {
			return nil, fmt.Errorf("failed to build query to insert chat_user: %v", err)
		}

		_, err = s.pool.Exec(ctx, sql, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to insert chat_user: %v", err)
		}

	}

	return &desc.CreateChatResponse{Id: chatID}, nil
}

func (s *server) DeleteChat(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {

	query := sq.Select("1").
		From("chat").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var exist int

	err = s.pool.QueryRow(ctx, sql, args...).Scan(&exist)
	if err != nil {
		if err.Error() == errNoRows {
			return nil, fmt.Errorf("chat: %d not founded", req.GetId())
		}
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	queryDelete := sq.Delete("chat").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": req.GetId()})

	sql, args, err = queryDelete.ToSql()
	_, err = s.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {

	user, err := s.grpcClient.GetUser(ctx, &descUser.GetUserRequest{
		Query: &descUser.GetUserRequest_Name{Name: req.GetFrom()},
	})
	if err != nil {
		return nil, err
	}

	// Пока выберу первый попавщийся, потом будет совсем иначе все
	querySelect := sq.Select("chat_id").
		From("chat_user").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"user_id": user.GetId()}).
		Limit(1)

	sql, args, err := querySelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}
	var chatID int64

	err = s.pool.QueryRow(ctx, sql, args...).Scan(&chatID)
	if err != nil {
		if err.Error() == errNoRows {
			return nil, fmt.Errorf("no chat found for user: %v", user.GetId())
		}
		return nil, fmt.Errorf("failed to select chat id: %v", err)
	}

	queryInsert := sq.Insert("message").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "user_id", "text").
		Values(chatID, user.GetId(), req.GetText())

	sql, args, err = queryInsert.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query to insert message: %v", err)
	}

	_, err = s.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert message: %v", err)
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

	conn, err := grpc.Dial(grpcClientConfig.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)

	client := descUser.NewAuthServiceClient(conn)

	desc.RegisterChatServerServiceServer(s, &server{pool: pool, grpcClient: client})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

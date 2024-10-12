package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/waryataw/chat-server/internal/config"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	descUser "github.com/waryataw/auth/pkg/user_v1"
	desc "github.com/waryataw/chat-server/pkg/chat_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatV1Server
	pool       *pgxpool.Pool
	clientConn *grpc.ClientConn
}

func (s *server) Create(ctx context.Context, _ *desc.CreateRequest) (*desc.CreateResponse, error) {
	client := descUser.NewUserV1Client(s.clientConn)

	user, err := client.GetByName(ctx, &descUser.GetByNameRequest{
		Name: "magna",
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(user)

	return &desc.CreateResponse{Id: gofakeit.Int64()}, nil
}

func (s *server) Delete(context.Context, *desc.DeleteRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(context.Context, *desc.SendMessageRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func main() {

	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
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
	defer conn.Close()

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
	desc.RegisterChatV1Server(s, &server{pool: pool, clientConn: conn})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package server

import (
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	chatTransport "github.com/sergripenko/chatRPC/internal/delivery/grpc"
	"github.com/sergripenko/chatRPC/internal/repository/mem"
	"github.com/sergripenko/chatRPC/internal/service/auth"
	"github.com/sergripenko/chatRPC/internal/service/chat"
	pb "github.com/sergripenko/chatRPC/protofiles"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() error {
	logFile, err := os.OpenFile("logs.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer logFile.Close()
	logrus.SetOutput(logFile)

	// create listener
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	repo := mem.NewInMemoryRepositoryService()
	authInterceptor := chatTransport.NewAuthInterceptor(auth.NewAuthService(repo))

	// create gRPC server
	server := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Unary()))
	chatService := chat.NewChatService(repo, repo, repo)
	grpcHandler := chatTransport.NewHandler(chatService)
	pb.RegisterChatServiceServer(server, grpcHandler)

	logrus.Info("start server")
	if err = server.Serve(listener); err != nil {
		return err
	}
	return err
}

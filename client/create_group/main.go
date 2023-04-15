package main

import (
	"context"
	"flag"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sergripenko/chatRPC/client"
	pb "github.com/sergripenko/chatRPC/protofiles"
)

var (
	username = flag.String("username", "", "Username for connect")
	group    = flag.String("group", "", "Group name")
)

func main() {
	flag.Parse()

	authInterceptor := client.NewAuthInterceptor(*username)
	// dial to server
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor.Unary()))

	if err != nil {
		logrus.Fatal("Error connecting to gRPC server: ", err)
	}
	defer conn.Close()

	// create the stream
	chatService := pb.NewChatServiceClient(conn)
	req := pb.CreateGroupChatRequest{GroupName: *group}
	ctx := context.Background()

	resp, err := chatService.CreateGroupChat(ctx, &req)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(resp)
}

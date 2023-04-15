package main

import (
	"context"
	"errors"
	"flag"
	"io"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/sergripenko/chatRPC/protofiles"
)

var (
	username = flag.String("username", "", "Username for connect")
)

func main() {
	flag.Parse()

	// dial to server
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logrus.Fatal("Error connecting to gRPC server: ", err)
	}

	defer conn.Close()

	// create the stream
	client := pb.NewChatServiceClient(conn)
	req := pb.ConnectRequest{Username: *username}

	stream, err := client.Connect(context.Background(), &req)
	if err != nil {
		logrus.Fatal(err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				logrus.Fatal(err)
			}
			continue
		}
		logrus.Info("sender: ", resp.Sender)
		logrus.Info("group: ", resp.Group)
		logrus.Info("message: ", resp.Message)
	}
}

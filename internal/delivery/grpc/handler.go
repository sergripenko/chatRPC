package grpc

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"

	"github.com/sergripenko/chatRPC/internal/domain"
	pb "github.com/sergripenko/chatRPC/protofiles"
)

type Handler struct {
	pb.UnimplementedChatServiceServer
	chatService ChatServiceProvider
}

func NewHandler(
	chatService ChatServiceProvider,
) *Handler {
	return &Handler{
		chatService: chatService,
	}
}

type ChatServiceProvider interface {
	Connect(ctx context.Context, user *domain.User) (<-chan *domain.Message, error)
	Disconnect(ctx context.Context, user *domain.User) error
	JoinGroupChat(ctx context.Context, user *domain.User, group *domain.Group) error
	LeaveGroupChat(ctx context.Context, user *domain.User, group *domain.Group) error
	CreateGroupChat(ctx context.Context, user *domain.User, group *domain.Group) error
	SendMessage(ctx context.Context, user *domain.User, mess *domain.Message) error
	ListChannels(ctx context.Context) ([]*domain.Group, error)
}

func (h *Handler) getUserFromCtx(ctx context.Context) (*domain.User, error) {
	md, _ := metadata.FromOutgoingContext(ctx)
	username := md["user"]
	if len(username) == 0 || strings.EqualFold(username[0], "") {
		return nil, errors.New("no user in ctx")
	}
	return &domain.User{Username: username[0]}, nil
}

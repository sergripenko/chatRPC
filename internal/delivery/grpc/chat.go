package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sergripenko/chatRPC/internal/domain"
	pb "github.com/sergripenko/chatRPC/protofiles"
)

func (h *Handler) Connect(req *pb.ConnectRequest, srv pb.ChatService_ConnectServer) error {
	username := req.GetUsername()
	if strings.EqualFold(username, "") {
		return status.Error(codes.InvalidArgument, "empty `username` param")
	}

	msgChan, err := h.chatService.Connect(srv.Context(), &domain.User{Username: req.GetUsername()})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for {
		select {
		case <-srv.Context().Done():
			return h.chatService.Disconnect(srv.Context(), &domain.User{Username: req.GetUsername()})
		case mess := <-msgChan:
			resp := pb.Message{
				Sender:  mess.Sender,
				Group:   mess.Group,
				Message: mess.Message,
			}
			if err = srv.Send(&resp); err != nil {
				return err
			}
		}
	}
}

func (h *Handler) JoinGroupChat(ctx context.Context, request *pb.JoinGroupChatRequest) (
	*pb.JoinGroupChatResponse, error) {
	user, err := h.getUserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if strings.EqualFold(request.GetGroupName(), "") {
		return nil, status.Error(codes.InvalidArgument, "invalid input data")
	}
	if err = h.chatService.JoinGroupChat(ctx, user, &domain.Group{Name: request.GroupName}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.JoinGroupChatResponse{}, nil
}

func (h *Handler) LeaveGroupChat(ctx context.Context, request *pb.LeaveGroupChatRequest) (
	*pb.LeaveGroupChatResponse, error) {
	user, err := h.getUserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if strings.EqualFold(request.GetGroupName(), "") {
		return nil, status.Error(codes.InvalidArgument, "invalid input data")
	}
	if err = h.chatService.LeaveGroupChat(ctx, user, &domain.Group{Name: request.GetGroupName()}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.LeaveGroupChatResponse{}, nil
}

func (h *Handler) CreateGroupChat(ctx context.Context, request *pb.CreateGroupChatRequest) (
	*pb.CreateGroupChatResponse, error) {
	user, err := h.getUserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if strings.EqualFold(request.GetGroupName(), "") {
		return nil, status.Error(codes.InvalidArgument, "invalid input data")
	}
	if err = h.chatService.CreateGroupChat(ctx, user, &domain.Group{Name: request.GetGroupName()}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateGroupChatResponse{}, nil
}

func (h *Handler) SendMessage(ctx context.Context, request *pb.SendMessageRequest) (
	*pb.SendMessageResponse, error) {
	user, err := h.getUserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	// validate msg type
	if (strings.EqualFold(request.GetUsername(), "") && strings.EqualFold(request.GetChannelName(), "")) ||
		(!strings.EqualFold(request.GetUsername(), "") && !strings.EqualFold(request.GetChannelName(), "")) {
		return nil, status.Error(codes.InvalidArgument, "invalid input data")
	}
	// validate msg input
	if strings.EqualFold(request.GetMessage(), "") {
		return nil, status.Error(codes.InvalidArgument, "empty message")
	}

	var msgType, receiver string
	if !strings.EqualFold(request.GetUsername(), "") {
		msgType = domain.DirectMessageType
		receiver = request.GetUsername()
	} else {
		msgType = domain.GroupMessageType
		receiver = request.GetChannelName()
	}

	mess := &domain.Message{
		Type:     msgType,
		Sender:   user.Username,
		Receiver: receiver,
		Group:    request.GetChannelName(),
		Message:  request.GetMessage(),
	}

	if err = h.chatService.SendMessage(ctx, &domain.User{Username: user.Username}, mess); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.SendMessageResponse{}, nil
}

func (h *Handler) ListChannels(ctx context.Context, request *pb.ListChannelsRequest) (
	*pb.ListChannelsResponse, error) {
	groupObjs, err := h.chatService.ListChannels(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var groups []*pb.Group

	for idx := range groupObjs {
		tmpGroup := &pb.Group{Name: groupObjs[idx].Name}

		for _, us := range groupObjs[idx].Users {
			tmpGroup.Users = append(tmpGroup.Users, us.Username)
		}
		groups = append(groups, tmpGroup)
	}
	return &pb.ListChannelsResponse{Group: groups}, nil
}

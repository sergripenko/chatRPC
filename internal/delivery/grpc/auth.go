package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/sergripenko/chatRPC/internal/domain"
)

type AuthInterceptor struct {
	authService AuthServiceProvider
}

func NewAuthInterceptor(
	authService AuthServiceProvider,
) *AuthInterceptor {
	return &AuthInterceptor{
		authService: authService,
	}
}

type AuthServiceProvider interface {
	Authorize(ctx context.Context, user *domain.User) (*domain.User, error)
}

func (ai *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		//logrus.Println("--> unary interceptor: ", info.FullMethod)

		user, err := ai.authorize(ctx)
		if err != nil {
			return nil, err
		}
		// set username
		md, _ := metadata.FromIncomingContext(ctx)
		md.Append("user", user.Username)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return handler(ctx, req)
	}
}

func (ai *AuthInterceptor) authorize(ctx context.Context) (*domain.User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	values := md["authorization"]

	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	user, err := ai.authService.Authorize(ctx, &domain.User{Username: values[0]})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication error")
	}
	return user, nil
}

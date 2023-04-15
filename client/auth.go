package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	token string
}

func NewAuthInterceptor(
	token string,
) *AuthInterceptor {
	return &AuthInterceptor{token: token}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		//logrus.Infof("--> unary interceptor: %s", method)
		return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
	}
}
func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.token)
}

package interceptors

import (
	"context"

	"gozero_init/pkg/xcode"

	"google.golang.org/grpc"
)

// Server 错误拦截器
func ServerErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		return resp, xcode.FromError(err).Err()
	}
}

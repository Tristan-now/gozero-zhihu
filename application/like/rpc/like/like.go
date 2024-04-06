// Code generated by goctl. DO NOT EDIT.
// Source: like.proto

package like

import (
	"context"

	"gozero_init/application/like/rpc/service"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	IsThumbupRequest  = service.IsThumbupRequest
	IsThumbupResponse = service.IsThumbupResponse
	ThumbupRequest    = service.ThumbupRequest
	ThumbupResponse   = service.ThumbupResponse
	UserThumbup       = service.UserThumbup

	Like interface {
		Thumbup(ctx context.Context, in *ThumbupRequest, opts ...grpc.CallOption) (*ThumbupResponse, error)
		IsThumbup(ctx context.Context, in *IsThumbupRequest, opts ...grpc.CallOption) (*IsThumbupResponse, error)
	}

	defaultLike struct {
		cli zrpc.Client
	}
)

func NewLike(cli zrpc.Client) Like {
	return &defaultLike{
		cli: cli,
	}
}

func (m *defaultLike) Thumbup(ctx context.Context, in *ThumbupRequest, opts ...grpc.CallOption) (*ThumbupResponse, error) {
	client := service.NewLikeClient(m.cli.Conn())
	return client.Thumbup(ctx, in, opts...)
}

func (m *defaultLike) IsThumbup(ctx context.Context, in *IsThumbupRequest, opts ...grpc.CallOption) (*IsThumbupResponse, error) {
	client := service.NewLikeClient(m.cli.Conn())
	return client.IsThumbup(ctx, in, opts...)
}

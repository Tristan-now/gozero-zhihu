package logic

import (
	"context"

	"gozero_init/application/user/rpc/internal/svc"
	"gozero_init/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIDLogic {
	return &FindByIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIDLogic) FindByID(in *service.FindByIdRequest) (*service.FindByIdResponse, error) {
	// todo: add your logic here and delete this line

	return &service.FindByIdResponse{}, nil
}

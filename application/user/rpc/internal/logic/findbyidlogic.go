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
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		logx.Errorf("FIndById userId: %d error: %v", in.UserId, err)
		return nil, err
	}

	return &service.FindByIdResponse{
		UserId:   user.Id,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}

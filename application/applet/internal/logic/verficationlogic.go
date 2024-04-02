package logic

import (
	"context"

	"gozero_init/application/applet/internal/svc"
	"gozero_init/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerficationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerficationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerficationLogic {
	return &VerficationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerficationLogic) Verfication(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	// todo: add your logic here and delete this line

	return
}

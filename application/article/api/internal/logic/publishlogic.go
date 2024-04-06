package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gozero_init/application/article/api/internal/code"
	"gozero_init/application/article/api/internal/svc"
	"gozero_init/application/article/api/internal/types"
	"gozero_init/application/article/rpc/pb"
)

const minContentLen = 80

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishRequest) (*types.PublishResponse, error) {
	if len(req.Title) == 0 {
		return nil, code.ArtitleTitleEmpty
	}
	if len(req.Content) < minContentLen {
		return nil, code.ArticleContentTooFewWords
	}
	if len(req.Cover) == 0 {
		return nil, code.ArticleCoverEmpty
	}
	//userId, err := l.ctx.Value("userId").(json.Number).Int64()
	userId := int64(2)
	//if err != nil {
	//	logx.Errorf("l.ctx.Value error: %v", err)
	//	return nil, xcode.NoLogin
	//}

	pret, err := l.svcCtx.ArticleRPC.Publish(l.ctx, &pb.PublishRequest{
		UserId:      userId,
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Cover:       req.Cover,
	})
	if err != nil {
		logx.Errorf("l.svcCtx.ArticleRPC.Publish req: %v userId: %d error: %v", req, userId, err)
		return nil, err
	}

	return &types.PublishResponse{ArticleId: pret.ArticleId}, nil
}

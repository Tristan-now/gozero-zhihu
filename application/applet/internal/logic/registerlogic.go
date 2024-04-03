package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gozero_init/application/applet/internal/code"
	"gozero_init/application/user/rpc/user"
	"gozero_init/pkg/encrypt"
	"gozero_init/pkg/jwt"
	"strings"

	"gozero_init/application/applet/internal/svc"
	"gozero_init/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	preFixActivation = "biz#activation#%s"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func checkVerificationCode(rds *redis.Redis, mobile, code string) error {
	cacheCode, err := getActivationCache(mobile, rds)
	if err != nil {
		return err
	}
	if cacheCode == "" {
		return errors.New("verification code expired")
	}
	if cacheCode != code {
		return errors.New("verification code failed")
	}
	return nil
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// todo: add your logic here and delete this line
	req.Name = strings.TrimSpace(req.Name)
	req.Mobile = strings.TrimSpace(req.Mobile)
	if len(req.Mobile) == 0 {
		return nil, code.RegisterMobileEmpty
	}
	req.Password = strings.TrimSpace(req.Password)
	if len(req.Password) == 0 {
		return nil, code.RegisterPasswdEmpty
	} else {
		req.Password = encrypt.EncPassword(req.Password)
	}

	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	if len(req.VerificationCode) == 0 {
		return nil, code.VerificationCodeEmpty
	}

	err = checkVerificationCode(l.svcCtx.BizRedis, req.Mobile, req.VerificationCode)
	if err != nil {
		logx.Errorf("checkVerificationCode error: %v", err)
		return nil, err
	}

	mobile, err := encrypt.EncMobile(req.Mobile)
	if err != nil {
		logx.Errorf("EncMobile mobileL %s error: %v", req.Mobile, err)
		return nil, err
	}

	u, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &user.FindByMobileRequest{
		Mobile: mobile,
	})
	if err != nil {
		logx.Errorf("FIndByMobile error: %v", err)
		return nil, err
	}

	if u != nil && u.UserId > 0 {
		return nil, code.MobileHasRegistered
	}
	regRet, err := l.svcCtx.UserRPC.Register(l.ctx, &user.RegisterRequest{
		Username: req.Name,
		Mobile:   mobile,
	})

	if err != nil {
		logx.Errorf("Register error: %v", err)
		return nil, err
	}

	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userID": regRet.UserId,
		},
	})

	if err != nil {
		logx.Errorf("BuildTokens error: %v", err)
		return nil, err
	}

	_ = delActivationCache(req.Mobile, req.VerificationCode, l.svcCtx.BizRedis)

	return &types.RegisterResponse{
		UserId: regRet.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}

package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gozero_init/application/user/rpc/user"
	"gozero_init/pkg/util"
	"strconv"
	"time"

	"gozero_init/application/applet/internal/svc"
	"gozero_init/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	prefixVerificationCount = "biz#verification#count#%s"
	verificationLimitPerDay = 10
	expireActivation        = 60 * 30
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

func (l *VerficationLogic) getVerificationCount(mobile string) (int, error) {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	val, err := l.svcCtx.BizRedis.Get(key)
	if err != nil {
		return 0, err
	}
	if len(val) == 0 {
		return 0, nil
	}

	return strconv.Atoi(val)
}

func getActivationCache(mobile string, rds *redis.Redis) (string, error) {
	key := fmt.Sprintf(preFixActivation, mobile)
	return rds.Get(key)
}

func saveActivationCache(mobile, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(preFixActivation, mobile)
	return rds.Setex(key, code, expireActivation)
}

func (l *VerficationLogic) incrVerificationCount(mobile string) error {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	_, err := l.svcCtx.BizRedis.Incr(key)
	if err != nil {
		return err
	}

	return l.svcCtx.BizRedis.Expireat(key, util.EndOfDay(time.Now()).Unix())
}

func (l *VerficationLogic) Verfication(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	// todo: add your logic here and delete this line
	count, err := l.getVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("getVerificationCount mobvile:%s error : %v", req.Mobile, err)
	}
	if count > verificationLimitPerDay {
		return nil, err
	}
	code, err := getActivationCache(req.Mobile, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("getActivationCache mobile: %s error: %v", req.Mobile, err)
	}
	if len(code) == 0 {
		code = util.RandomNumeric(6)
	}

	_, err = l.svcCtx.UserRPC.SendSms(l.ctx, &user.SendSmsRequest{
		Mobile: req.Mobile,
	})

	if err != nil {
		logx.Errorf("sendSms mobile: %s error: %v", req.Mobile, err)
	}

	err = saveActivationCache(req.Mobile, code, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("saveActivationCache mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}

	err = l.incrVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("incrVerificationCount mobile: %s error: %v", req.Mobile, err)
	}

	return &types.VerificationResponse{}, nil
}

func delActivationCache(mobile, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(preFixActivation, mobile)
	_, err := rds.Del(key)
	return err
}

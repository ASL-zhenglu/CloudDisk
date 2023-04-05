package logic

import (
	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/models"
	"context"
	"errors"
	"fmt"
	"time"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MailCodeSendRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMailCodeSendRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MailCodeSendRegisterLogic {
	return &MailCodeSendRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MailCodeSendRegisterLogic) MailCodeSendRegister(req *types.MailCodeSendRequest) (resp *types.MailCodeSendReply, err error) {
	// todo: add your logic here and delete this line
	// 该邮箱未被注册

	cnt, err := l.svcCtx.Engine.Where("email = ?", req.Email).Count(new(models.UserBasic))

	//cnt, err := l.svcCtx.Engine.Where("email = ?", req.Email).Count(new(models.UserBasic))
	if err != nil {
		return
	}
	if cnt > 0 {
		err = errors.New("该邮箱已经注册")
		return
	}
	fmt.Println(cnt, req.Email)
	fmt.Println("?????????????????")
	// 生成随机数
	code := helper.RandCode()
	//存储验证码
	l.svcCtx.RDB.Set(l.ctx, req.Email, code, time.Second*time.Duration(define.CodeExpire))
	// 发送验证码
	//go func() {
	//	err = helper.MailSendCode(req.Email, code)
	//}()
	err = helper.MailSendCode(req.Email, code)
	return
}

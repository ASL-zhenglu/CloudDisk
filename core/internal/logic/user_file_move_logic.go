package logic

import (
	"cloud-disk/core/models"
	"context"
	"errors"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileMoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileMoveLogic {
	return &UserFileMoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileMoveLogic) UserFileMove(req *types.UserFileMoveRequest, userIdentity string) (resp *types.UserFileMoveReply, err error) {
	// todo: add your logic here and delete this line
	// 查ParentId
	parentDate := new(models.UserRepository)
	has, err := l.svcCtx.Engine.Where("identity = ? AND user_identity = ?", req.ParentIdentity, userIdentity).Get(parentDate)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("文件夹不存在")
	}
	// 更新记录的parentId
	_, err = l.svcCtx.Engine.Where("identity = ?", req.Identity).Update(models.UserRepository{
		ParentId: int64(parentDate.Id),
	})
	return
}

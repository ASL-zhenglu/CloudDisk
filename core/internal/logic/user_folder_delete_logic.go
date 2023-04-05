package logic

import (
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFolderDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFolderDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFolderDeleteLogic {
	return &UserFolderDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFolderDeleteLogic) UserFolderDelete(req *types.UserFolderDeleteRequest, userIdentity string) (resp *types.UserFolderDeleteReply, err error) {
	// todo: add your logic here and delete this line
	_, err = l.svcCtx.Engine.Where("user_identity = ? AND identity = ?", userIdentity, req.Identity).Delete(new(models.UserRepository))
	if err != nil {
		return
	}
	return
}

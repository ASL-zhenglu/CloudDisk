package logic

import (
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadPrepareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadPrepareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadPrepareLogic {
	return &FileUploadPrepareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadPrepareLogic) FileUploadPrepare(req *types.FileUploadPrepareRequest) (resp *types.FileUploadPrepareReply, err error) {
	rp := new(models.RepositoryPool)
	has, err := l.svcCtx.Engine.Where("hash=?", req.Md5).Get(rp)
	if err != nil {
		return
	}
	resp = new(types.FileUploadPrepareReply)
	if has {
		// 秒传
		resp.Identity = rp.Identity
	} else {
		// 分片上传
		key, uploadId, err := helper.CosInitPart(req.Ext)
		if err != nil {
			return nil, err
		}
		resp.Key = key
		resp.UploadId = uploadId
	}
	return
}
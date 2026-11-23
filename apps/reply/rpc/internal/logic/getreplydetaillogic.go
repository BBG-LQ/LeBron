package logic

import (
	"context"

	"lebron/apps/reply/rpc/internal/svc"
	"lebron/apps/reply/rpc/reply"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReplyDetailLogic {
	return &GetReplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取评论详情
func (l *GetReplyDetailLogic) GetReplyDetail(in *reply.GetReplyDetailReq) (*reply.GetReplyDetailResp, error) {
	// todo: add your logic here and delete this line

	re, err := l.svcCtx.ReplyModel.FindOne2(l.ctx, uint64(in.ReplyId))
	if err != nil {
		return nil, err
	}

	data := &reply.ReplyInfo{
		TargetId:   int64(re.TargetId),
		TargetType: int64(re.TargetType),
		ParentId:   int64(re.ParentId),
		Content:    re.Content,
	}

	return &reply.GetReplyDetailResp{
		Reply: data,
	}, nil
}

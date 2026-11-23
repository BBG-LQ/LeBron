package logic

import (
	"context"

	"lebron/apps/reply/rpc/internal/svc"
	"lebron/apps/reply/rpc/reply"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeReplyLogic {
	return &LikeReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 点赞/取消点赞评论
func (l *LikeReplyLogic) LikeReply(in *reply.LikeReplyReq) (*reply.LikeReplyResp, error) {
	// todo: add your logic here and delete this line

	var num uint64

	re, err := l.svcCtx.ReplyModel.FindOne2(l.ctx, uint64(in.ReplyId))
	if err != nil {
		return nil, err
	}

	num = re.LikeCount
	if in.Like {
		num = re.LikeCount + 1
	}

	err = l.svcCtx.ReplyModel.Like(l.ctx, num, uint64(in.ReplyId))
	if err != nil {
		return nil, err
	}
	succ := true
	mess := "no errr"
	return &reply.LikeReplyResp{
		Success:      succ,
		Message:      mess,
		NewLikeCount: int64(num),
	}, nil
}

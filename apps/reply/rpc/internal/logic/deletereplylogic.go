package logic

import (
	"context"

	"lebron/apps/reply/rpc/internal/svc"
	"lebron/apps/reply/rpc/reply"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteReplyLogic {
	return &DeleteReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除评论（软删除或硬删除）
func (l *DeleteReplyLogic) DeleteReply(in *reply.DeleteReplyReq) (*reply.DeleteReplyResp, error) {

	// todo: add your logic here and delete this line

	var succ bool
	var mess string

	replyDetail, err := l.svcCtx.ReplyModel.FindOne(l.ctx, uint64(in.ReplyId))
	if err != nil {
		mess = "not fount"
		return nil, err
	}

	if replyDetail.UserId == uint64(in.UserId) {
		err = l.svcCtx.ReplyModel.DeleteWithuidandRid(l.ctx, uint64(in.ReplyId), uint64(in.UserId))
		if err != nil {
			mess = "not uid"
			return nil, err

		}
		succ = true
	}

	return &reply.DeleteReplyResp{
		Success: succ,
		Message: mess,
	}, nil
}

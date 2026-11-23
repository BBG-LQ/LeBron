package logic

import (
	"context"
	"encoding/json"

	"lebron/apps/reply/rpc/internal/model"
	"lebron/apps/reply/rpc/internal/svc"
	"lebron/apps/reply/rpc/reply"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReplyLogic {
	return &UpdateReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新评论内容（仅限作者本人）
func (l *UpdateReplyLogic) UpdateReply(in *reply.UpdateReplyReq) (*reply.UpdateReplyResp, error) {
	// todo: add your logic here and delete this line

	rep, err := l.svcCtx.ReplyModel.FindOne(l.ctx, uint64(in.ReplyId))
	if err != nil {
		return nil, err
	}

	if rep.UserId != uint64(in.UserId) {
		return nil, err
	}

	imagesJSON, _ := json.Marshal(in.Images)
	data := &model.Reply{
		ReplyId: uint64(in.ReplyId),
		UserId:  uint64(in.UserId),
		Content: in.Content,
		Images:  string(imagesJSON),
		Video:   in.Video,
	}

	l.svcCtx.ReplyModel.Update(l.ctx, data)

	return &reply.UpdateReplyResp{
		Success: true,
		Message: "111",
	}, nil
}

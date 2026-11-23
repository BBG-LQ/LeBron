package logic

import (
	"context"

	"lebron/apps/reply/rpc/internal/svc"
	"lebron/apps/reply/rpc/reply"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportReplyLogic {
	return &ReportReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 举报评论
func (l *ReportReplyLogic) ReportReply(in *reply.ReportReplyReq) (*reply.ReportReplyResp, error) {
	// todo: add your logic here and delete this line

	return &reply.ReportReplyResp{}, nil
}

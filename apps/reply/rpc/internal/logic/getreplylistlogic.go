package logic

import (
	"context"
	"time"

	"lebron/apps/reply/rpc/internal/model"
	"lebron/apps/reply/rpc/internal/svc"
	"lebron/apps/reply/rpc/reply"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReplyListLogic {
	return &GetReplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取评论列表（支持按目标对象分页，可包含子回复）
func (l *GetReplyListLogic) GetReplyList(in *reply.GetReplyListReq) (*reply.GetReplyListResp, error) {
	// todo: add your logic here and delete this line

	if in.Cursor == 0 {
		in.Cursor = time.Now().Unix()
	}
	if in.Ps == 0 {
		in.Ps = 2
	}

	var (
		isEnd     bool
		lastTime  int64
		firstPage []*reply.ReplyItem
		replys    []*model.Reply
	)

	var (
		err   error
		ctime = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")
	)
	//读数据库

	replys, err = l.svcCtx.ReplyModel.ReplyList(l.ctx, uint64(in.TargetType), uint64(in.TargetId), uint64(in.ParentId), ctime, int64(in.Ps), in.Order)
	if err != nil {
		return nil, err
	}
	//
	var firstPageProducts []*model.Reply
	if len(replys) > int(in.Ps) {
		firstPageProducts = replys[:int(in.Ps)]

	} else {
		firstPageProducts = replys
		isEnd = true
	}

	//
	count := len(replys)
	for _, p := range firstPageProducts {
		rep := &reply.ReplyInfo{
			ReplyId:    int64(p.ReplyId),
			UserId:     int64(p.UserId),
			TargetId:   int64(p.TargetId),
			TargetType: int64(p.TargetType),
			Content:    p.Content,
			ReplyCount: int64(count),
		}

		firstPage = append(firstPage, &reply.ReplyItem{
			Reply: rep,
		})
	}

	//

	if len(firstPage) > 0 {
		pageLast := firstPage[len(firstPage)-1]

		lastTime = pageLast.Reply.CreateTime

		if lastTime < 0 {
			lastTime = 0
		}

		// 去重处理：如果客户端传入了 productId（用于断点续传），则跳过之前已返回的记录
		// 这里假设前端会传入上次返回的最后一条记录的 productId 和 cursor，用于避免重复

		for k, p := range firstPage {
			if p.Reply.CreateTime == in.Cursor {
				firstPage = firstPage[k:] // 截取从匹配项开始
				break
			}
		}

	}

	ret := &reply.GetReplyListResp{
		Items:      firstPage,
		NextCursor: lastTime,
		IsEnd:      isEnd,
		TotalCount: int64(count),
	}

	return ret, nil
}

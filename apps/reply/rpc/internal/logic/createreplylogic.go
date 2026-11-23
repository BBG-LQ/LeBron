package logic

import (
	"context"
	"encoding/json"

	"lebron/apps/reply/rpc/internal/model"
	"lebron/apps/reply/rpc/internal/svc"
	"lebron/apps/reply/rpc/reply"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReplyLogic {
	return &CreateReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建评论（或回复）
func (l *CreateReplyLogic) CreateReply(in *reply.CreateReplyReq) (*reply.CreateReplyResp, error) {
	// todo: add your logic here and delete this line

	//create 新的评论
	imagesJSON, _ := json.Marshal(in.Images)
	var succ bool
	var mess string

	var num1, num2 uint64

	data := &model.Reply{

		UserId:     uint64(in.UserId),
		TargetId:   uint64(in.TargetId),
		TargetType: uint64(in.TargetType),
		ParentId:   uint64(in.ParentId),
		Content:    in.Content,
		Images:     string(imagesJSON),
		Video:      in.GetVideo(),
	}

	result, err := l.svcCtx.ReplyModel.Insert(l.ctx, data)
	if err != nil {
		mess = "not insert"
		return nil, err
	}
	_, err = result.LastInsertId()

	if err != nil {
		mess = "not insert"
		return nil, err
	}

	succ = true

	// count

	if in.ParentId > 0 {

		re, err := l.svcCtx.ReplyModel.FindOne2(l.ctx, uint64(in.ParentId))
		if err != nil {
			return nil, err
		}
		mess = "right"
		num1 = re.ReplyId
		num2 = re.ReplyCount

		if re.ParentId == 0 {

			mess = re.Content
			var count uint64

			count = num2 + 1
			err := l.svcCtx.ReplyModel.UpdateRECount(l.ctx, uint64(in.ParentId), count)
			if err != nil {
				return nil, err
			}

		}
	}

	// succ = true
	return &reply.CreateReplyResp{

		ReplyId: int64(num1),
		Success: succ,
		Message: mess,
		Num:     int64(num2),
	}, nil
}

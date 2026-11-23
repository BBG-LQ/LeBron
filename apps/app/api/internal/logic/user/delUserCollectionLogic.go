// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"lebron/apps/app/api/internal/svc"
	"lebron/apps/app/api/internal/types"
	"lebron/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelUserCollectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// del user collection
func NewDelUserCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserCollectionLogic {
	return &DelUserCollectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelUserCollectionLogic) DelUserCollection(req *types.UserCollectionDelReq) (resp *types.UserCollectionDelRes, err error) {
	// todo: add your logic here and delete this line

	var delCollectRpcRep user.DeleteFavoriteRequest
	delCollectRpcRep.Id = req.Id
	_, err = l.svcCtx.UserRPC.DeleteFavorite(l.ctx, &delCollectRpcRep)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

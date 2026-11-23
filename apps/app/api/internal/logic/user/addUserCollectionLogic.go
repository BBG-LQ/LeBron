// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"encoding/json"

	"lebron/apps/app/api/internal/svc"
	"lebron/apps/app/api/internal/types"
	"lebron/apps/user/rpc/user"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserCollectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// add user collection
func NewAddUserCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserCollectionLogic {
	return &AddUserCollectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddUserCollectionLogic) AddUserCollection(req *types.UserCollectionAddReq) (resp *types.UserCollectionAddRes, err error) {
	// todo: add your logic here and delete this line
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	var addCollectionRpcReq user.AddFavoriteRequest
	addCollectionRpcReq.Uid = uid

	cppyErr := copier.Copy(&addCollectionRpcReq, req)
	if cppyErr != nil {
		return nil, cppyErr
	}

	_, err = l.svcCtx.UserRPC.AddFavorite(l.ctx, &addCollectionRpcReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

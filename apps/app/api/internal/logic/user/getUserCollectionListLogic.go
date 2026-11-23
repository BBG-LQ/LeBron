// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"encoding/json"

	"lebron/apps/app/api/internal/svc"
	"lebron/apps/app/api/internal/types"
	"lebron/apps/user/rpc/user"
	"lebron/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserCollectionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get user collection list
func NewGetUserCollectionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCollectionListLogic {
	return &GetUserCollectionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserCollectionListLogic) GetUserCollectionList(req *types.UserCollectionListReq) (resp *types.UserCollectionListRes, err error) {
	// todo: add your logic here and delete this line
	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("Error! get uid from token"), "Failed toget uid from token err : %v ,req:%+v", err, req)
	}

	var getCollectListReq user.ListFavoritesRequest

	getCollectListReq.Uid = uid

	var productIDs []int64
	products, err := l.svcCtx.UserRPC.ListFavorites(l.ctx, &getCollectListReq)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("Error! Function UserCollectionList"), "Failed to get user UserCollectionList  err : %v ,req:%+v", err, req)
	}

	for _, productRpclist := range products.ProductId {

		productIDs = append(productIDs, productRpclist)

	}

	return &types.UserCollectionListRes{
		ProductId: productIDs,
	}, nil
}

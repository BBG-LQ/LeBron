package logic

import (
	"context"

	"lebron/apps/user/rpc/internal/svc"
	"lebron/apps/user/rpc/model"
	"lebron/apps/user/rpc/user"
	"lebron/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFavoriteLogic {
	return &AddFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 添加收藏
func (l *AddFavoriteLogic) AddFavorite(in *user.AddFavoriteRequest) (*user.AddFavoriteResponse, error) {
	// todo: add your logic here and delete this line
	dbCollection := new(model.UserCollection)

	dbCollection.Uid = uint64(in.GetUid())
	dbCollection.ProductId = uint64(in.GetProductId())

	_, err := l.svcCtx.UserCollectionModel.Insert(l.ctx, dbCollection)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "AddUserCollection Database Exception : %+v , err: %v", dbCollection, err)
	}

	return &user.AddFavoriteResponse{}, nil
}

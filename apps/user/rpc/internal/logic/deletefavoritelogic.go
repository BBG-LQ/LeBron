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

type DeleteFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFavoriteLogic {
	return &DeleteFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除收藏
func (l *DeleteFavoriteLogic) DeleteFavorite(in *user.DeleteFavoriteRequest) (*user.DeleteFavoriteResponse, error) {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.UserCollectionModel.FindOne(l.ctx, uint64(in.GetId()))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.Wrap(xerr.NewErrMsg("数据不存在"), "该商品没有被收藏")
		}
		return nil, err
	}

	dbCollection := new(model.UserCollection)
	dbCollection.Id = uint64(in.GetId())
	dbCollection.IsDelete = 1
	err = l.svcCtx.UserCollectionModel.UpdateIsDelete(l.ctx, dbCollection)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "DelUserCollection Database Exception : %+v , err: %v", dbCollection, err)
	}

	return &user.DeleteFavoriteResponse{}, nil
}

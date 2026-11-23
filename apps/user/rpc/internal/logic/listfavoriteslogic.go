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

type ListFavoritesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFavoritesLogic {
	return &ListFavoritesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取收藏列表
func (l *ListFavoritesLogic) ListFavorites(in *user.ListFavoritesRequest) (*user.ListFavoritesResponse, error) {
	// todo: add your logic here and delete this line
	listFavrits, err := l.svcCtx.UserCollectionModel.FindAllByUid(l.ctx, uint64(in.Uid))

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "Failed  get user's Collection list err : %v , in :%+v", err, in)
	}

	var listPro []int64
	for _, v := range listFavrits {
		listPro = append(listPro, int64(v.ProductId))
	}

	return &user.ListFavoritesResponse{
		ProductId: listPro,
	}, nil
}

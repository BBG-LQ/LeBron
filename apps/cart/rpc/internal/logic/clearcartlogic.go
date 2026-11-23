package logic

import (
	"context"

	"lebron/apps/cart/rpc/cart"
	"lebron/apps/cart/rpc/internal/model"
	"lebron/apps/cart/rpc/internal/svc"
	"lebron/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ClearCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearCartLogic {
	return &ClearCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 清空购物车（删除该用户所有商品）
func (l *ClearCartLogic) ClearCart(in *cart.ClearCartReq) (*cart.ClearCartResp, error) {
	// todo: add your logic here and delete this line

	var success bool
	listCarts, err := l.svcCtx.CartModel.FindAllByUid(l.ctx, uint64(in.UserId))

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "Failed  get user's cart list err : %v , in :%+v", err, in)
	}

	for _, v := range listCarts {

		err = l.svcCtx.CartModel.Delete(l.ctx, v.Id)
		if err != nil {
			return nil, err
		}
	}

	success = true

	return &cart.ClearCartResp{
		Success: success,
		Message: "no",
	}, nil
}

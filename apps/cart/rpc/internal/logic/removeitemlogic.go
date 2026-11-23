package logic

import (
	"context"

	"lebron/apps/cart/rpc/cart"
	"lebron/apps/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveItemLogic {
	return &RemoveItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 从购物车中删除商品
func (l *RemoveItemLogic) RemoveItem(in *cart.RemoveItemReq) (*cart.RemoveItemResp, error) {
	// todo: add your logic here and delete this line

	cartPro, err := l.svcCtx.CartModel.FindOneByPid(l.ctx, uint64(in.ProductId))

	if err != nil {
		return nil, err
	}

	err = l.svcCtx.CartModel.Delete(l.ctx, cartPro.Id)
	if err != nil {
		return nil, err
	}

	return &cart.RemoveItemResp{
		Success: true,
		Message: "no",
	}, nil
}

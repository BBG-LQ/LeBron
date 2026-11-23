package logic

import (
	"context"
	"errors"

	"lebron/apps/cart/rpc/cart"
	"lebron/apps/cart/rpc/internal/model"
	"lebron/apps/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateItemLogic {
	return &UpdateItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新购物车中商品的数量
func (l *UpdateItemLogic) UpdateItem(in *cart.UpdateItemReq) (*cart.UpdateItemResp, error) {
	// todo: add your logic here and delete this line

	data := &model.Cart{
		Userid:    uint64(in.UserId),
		Productid: uint64(in.ProductId),
	}
	var succ bool
	var message string

	pods, err := l.svcCtx.CartModel.FindAllByUid(l.ctx, uint64(in.UserId))

	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, err
	}

	found := false

	for _, v := range pods {

		if v.Productid == uint64(in.ProductId) {

			data.Id = v.Id
			cartProQuan := v.Quantity
			data.Quantity = cartProQuan + in.Quantity

			err = l.svcCtx.CartModel.Update(l.ctx, data)

			if err != nil {
				return nil, err
			}
			found = true
			succ = true
			break

		}
	}

	if !found {
		message = "not fount"
	} else {
		message = "fount"
	}

	return &cart.UpdateItemResp{
		Success: succ,
		Message: message,
	}, nil
}

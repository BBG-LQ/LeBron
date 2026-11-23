package logic

import (
	"context"
	"errors"

	"lebron/apps/cart/rpc/cart"
	"lebron/apps/cart/rpc/internal/model"
	"lebron/apps/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddItemLogic {
	return &AddItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 添加商品到购物车（若商品已存在则增加数量）
func (l *AddItemLogic) AddItem(in *cart.AddItemReq) (*cart.AddItemResp, error) {
	// todo: add your logic here and delete this line

	data := &model.Cart{
		Userid:    uint64(in.UserId),
		Productid: uint64(in.ProductId),
	}

	var succ bool

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
		data.Quantity = in.Quantity
		_, err = l.svcCtx.CartModel.Insert(l.ctx, data)

		if err != nil {
			return nil, err
		}
		succ = true
	}

	// if err != nil {
	// 	if errors.Is(err, model.ErrNotFound) {

	// 		data.Quantity = in.Quantity
	// 		_, err = l.svcCtx.CartModel.Insert(l.ctx, data)

	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		succ = true

	// 	} else {
	// 		return nil, err
	// 	}

	// } else {
	// 	// for _, v := range pods {

	// 	// 	if v.Productid == uint64(in.ProductId) {

	// 	// 		data.Id = v.Id
	// 	// 		cartProQuan := v.Quantity
	// 	// 		data.Quantity = cartProQuan + in.Quantity

	// 	// 		err = l.svcCtx.CartModel.Update(l.ctx, data)

	// 	// 		if err != nil {
	// 	// 			return nil, err
	// 	// 		}
	// 	// 		succ = true

	// 	// 	}
	// 	// }
	// }

	return &cart.AddItemResp{
		Success: succ,
		Message: "no",
	}, nil
}

package logic

import (
	"context"
	"errors"

	"lebron/apps/cart/rpc/cart"
	"lebron/apps/cart/rpc/internal/model"
	"lebron/apps/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckItemLogic {
	return &CheckItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 选中/取消选中购物车中的商品（用于下单结算）
func (l *CheckItemLogic) CheckItem(in *cart.CheckItemReq) (*cart.CheckItemResp, error) {
	// todo: add your logic here and delete this line

	var checked int64

	if in.Checked {
		checked = 1
	} else {
		checked = 0
	}
	data := &model.Cart{
		Userid:    uint64(in.UserId),
		Productid: uint64(in.ProductId),

		Checked: checked,
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
			data.Quantity = v.Quantity

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

	return &cart.CheckItemResp{
		Success: succ,
		Message: message,
	}, nil
}

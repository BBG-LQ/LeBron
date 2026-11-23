package logic

import (
	"context"

	"lebron/apps/cart/rpc/cart"
	"lebron/apps/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCheckedItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger

	GetCartLogic *GetCartLogic
}

func NewGetCheckedItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCheckedItemsLogic {
	return &GetCheckedItemsLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		Logger:       logx.WithContext(ctx),
		GetCartLogic: NewGetCartLogic(ctx, svcCtx),
	}
}

// 批量获取购物车中选中商品的信息（用于订单创建前校验）
func (l *GetCheckedItemsLogic) GetCheckedItems(in *cart.GetCheckedItemsReq) (*cart.GetCheckedItemsResp, error) {
	// todo: add your logic here and delete this line

	getcartres := &cart.GetCartReq{
		UserId: in.UserId,
	}
	pros, err := l.GetCartLogic.GetCart(getcartres)

	if err != nil {
		return nil, err
	}

	var products []*cart.CartItem
	var TotalPrice float64

	for _, v := range pros.Cart.Items {

		if v.Checked {
			item := &cart.CartItem{

				ProductId:   v.ProductId,
				ProductName: v.ProductName,
				Price:       v.Price,
				Quantity:    v.Quantity,
				Checked:     v.Checked,
			}

			TotalPrice = TotalPrice + float64(v.Price)*float64(v.Quantity)

			products = append(products, item)
		}

	}

	return &cart.GetCheckedItemsResp{
		Items:      products,
		TotalPrice: TotalPrice,
	}, nil
}

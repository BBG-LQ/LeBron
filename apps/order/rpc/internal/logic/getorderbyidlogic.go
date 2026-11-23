package logic

import (
	"context"

	"lebron/apps/order/rpc/internal/svc"
	"lebron/apps/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderByIdLogic {
	return &GetOrderByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据订单ID查询订单详情（返回完整订单信息）
func (l *GetOrderByIdLogic) GetOrderById(in *order.GetOrderByIdReq) (*order.GetOrderByIdResp, error) {
	// todo: add your logic here and delete this line

	ord, err := l.svcCtx.OrdersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &order.GetOrderByIdResp{
		Order: &order.Orders{
			Userid:      int64(ord.Userid),
			Shoppingid:  ord.Shoppingid,
			Payment:     ord.Payment,
			Paymenttype: ord.Paymenttype,
		},
	}, nil
}

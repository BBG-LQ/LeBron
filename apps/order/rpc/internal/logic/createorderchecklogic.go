package logic

import (
	"context"
	"fmt"

	"lebron/apps/order/rpc/internal/svc"
	"lebron/apps/order/rpc/order"
	"lebron/apps/product/rpc/productclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderCheckLogic {
	return &CreateOrderCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建订单前的检查（如库存、价格等）
func (l *CreateOrderCheckLogic) CreateOrderCheck(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	// todo: add your logic here and delete this line

	fmt.Printf("CreateOrderCheck in:%v\n", in)
	res := productclient.ProductReq{
		ProductId: in.Pid,
	}
	pro, err := l.svcCtx.ProductRpc.GetProductDetail(l.ctx, &res)
	if err != nil {
		return nil, err
	}

	stock := pro.Stock
	price := pro.Price

	fmt.Printf("CreateOrderCheck stoc:%d\n", stock)
	fmt.Printf("CreateOrderCheck price:%f\n", price)

	return &order.CreateOrderResponse{}, nil
}

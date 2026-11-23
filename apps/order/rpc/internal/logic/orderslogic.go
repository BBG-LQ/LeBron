package logic

import (
	"context"
	"time"

	"lebron/apps/order/rpc/internal/svc"
	"lebron/apps/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrdersLogic {
	return &OrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrdersLogic) Orders(in *order.OrdersRequest) (*order.OrdersResponse, error) {
	// todo: add your logic here and delete this line
	ctime := time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")

	ords, err := l.svcCtx.OrderitemModel.FindByUid(l.ctx, in.UserId, ctime, int64(in.Ps))

	if err != nil {
		return nil, err
	}

	var IsEnd bool
	var ors, resp []*order.Orderitem

	for _, v := range ords {

		item := order.Orderitem{
			Id:               int64(v.Id),        //
			Orderid:          v.OrderId,          //
			Userid:           int64(v.UserId),    //
			Proid:            int64(v.ProductId), //
			Proname:          v.ProductName,      //
			Proimage:         v.ProductImage,     //
			Currentunitprice: v.CurrentPrice,     //
			Quantity:         v.Quantity,         //                          `
			Totalprice:       v.TotalPrice,       //
			CreateTime:       time.Now().Unix(),  //

		}

		ors = append(ors, &item)
	}

	if len(ors) > int(in.Ps) {
		resp = ors[:int(in.Ps)]

	} else {
		resp = ors
		IsEnd = true
	}

	return &order.OrdersResponse{
		Orders: resp,
		IsEnd:  IsEnd,
	}, nil
}

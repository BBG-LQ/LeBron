package logic

import (
	"context"

	"lebron/apps/order/rpc/internal/svc"
	"lebron/apps/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type RollbackOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRollbackOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RollbackOrderLogic {
	return &RollbackOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 回滚订单（用于事务补偿，例如分布式事务失败时撤销）
func (l *RollbackOrderLogic) RollbackOrder(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	// todo: add your logic here and delete this line

	return &order.CreateOrderResponse{}, nil
}

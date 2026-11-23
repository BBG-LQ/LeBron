package logic

import (
	"context"
	"database/sql"
	"fmt"

	"lebron/apps/order/rpc/internal/svc"
	"lebron/apps/order/rpc/order"
	"lebron/apps/user/rpc/userclient"

	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
)

type CreateOrderDTMRevertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderDTMRevertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderDTMRevertLogic {
	return &CreateOrderDTMRevertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DTM 反向操作：回滚创建订单（用于 Saga/TCC 模式的补偿）
func (l *CreateOrderDTMRevertLogic) CreateOrderDTMRevert(in *order.AddOrderReq) (*order.AddOrderResp, error) {
	// todo: add your logic here and delete this line
	//databasde
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	// 从 gRPC 上下文中获取 DTM 子事务屏障对象（用于处理分布式事务的幂等和补偿）
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {

		// 调用用户RPC服务，检查用户是否存在（若不存在则补偿操作失败，回滚事务）
		_, err := l.svcCtx.UserRpc.GetUser(l.ctx, &userclient.GetUserRequest{
			Id: in.Userid,
		})

		if err != nil {
			// 用户不存在时返回错误，触发事务回滚
			return fmt.Errorf("用户不存在")
		}
		//order

		resOrder, err := l.svcCtx.OrdersModel.FindOneByUid(l.ctx, in.Userid)
		if err != nil {
			// 如果查询不到订单，返回错误（可能订单尚未创建或已被删除）
			return fmt.Errorf("订单不存在")
		}

		// 将订单状态修改为60，表示订单已关闭（补偿操作）
		resOrder.Status = 60

		err = l.svcCtx.OrdersModel.TxUpdate(tx, resOrder.Id, resOrder.Status)
		if err != nil {
			// 更新失败则返回错误，回滚事务
			return fmt.Errorf("订单更新失败")
		}

		return nil

	}); err != nil {
		// 如果事务执行失败，包装错误并返回
		// 如果屏障执行或事务中出现错误，返回500错误
		return nil, status.Error(500, err.Error())
	}

	return &order.AddOrderResp{}, nil
}

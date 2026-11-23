package logic

import (
	"context"
	"database/sql"

	"lebron/apps/order/rpc/internal/svc"
	"lebron/apps/order/rpc/order"
	"lebron/pkg/xerr"

	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
)

type CancelCreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelCreateOrderLogic {
	return &CancelCreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelCreateOrderLogic) CancelCreateOrder(in *order.CancelCreateOrderReq) (*order.CancelCreateOrderResp, error) {
	// todo: add your logic here and delete this line

	// ===== 第一步：获取 DTM 子事务屏障 =====
	// 从 gRPC 上下文中获取 DTM 子事务屏障对象（用于处理分布式事务的幂等和补偿）
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(502, err.Error())
	}
	// ===== 第二步：获取数据库连接 =====
	//databasde
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(501, err.Error())
	}

	err = barrier.CallWithDB(db, func(tx *sql.Tx) error {

		// err := l.svcCtx.ShippingModel.TxUpdate(tx, uint64(in.ShippingId), 1)
		// if err != nil {
		// 	return err
		// }

		err = l.svcCtx.OrderitemModel.TxUpdate(tx, in.OrderId, 2)
		if err != nil {
			return err
		}

		err = l.svcCtx.OrdersModel.TxUpdate(tx, in.OrderId, 2)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "ConfirmCreateOrder failed: %v", err)
	}
	return &order.CancelCreateOrderResp{}, nil
}

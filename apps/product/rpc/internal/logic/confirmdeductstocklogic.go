package logic

import (
	"context"
	"database/sql"

	"lebron/apps/product/rpc/internal/svc"
	"lebron/apps/product/rpc/product"
	"lebron/pkg/xerr"

	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
)

type ConfirmDeductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmDeductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmDeductStockLogic {
	return &ConfirmDeductStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfirmDeductStockLogic) ConfirmDeductStock(in *product.ConfirmDeductStockReq) (*product.ConfirmDeductStockResp, error) {
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

	// ===== 第三步：在屏障中执行 Try 逻辑 =====

	// 开启子事务屏障（屏障会保证在数据库事务中执行，并且自动处理重复请求的幂等性）
	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {

		// 3. 冻结库存：将可用库存减少，冻结库存增加
		// 假设 product 表有 frozen_stock 字段

		_, err = tx.Exec(
			"UPDATE product SET status=555 WHERE id = ?",
			in.ProductId,
		)

		return nil
	}); err != nil {
		// 如果事务执行失败，包装错误并返回
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "product SET stock Exception : err: %v", err)
	}

	return &product.ConfirmDeductStockResp{}, nil
}

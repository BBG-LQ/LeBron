package logic

import (
	"context"
	"database/sql"

	"lebron/apps/order/rpc/internal/model"
	"lebron/apps/order/rpc/internal/svc"
	"lebron/apps/order/rpc/order"
	"lebron/pkg/xerr"

	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
)

type TryCreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTryCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TryCreateOrderLogic {
	return &TryCreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== TCC 参与者接口 =====
func (l *TryCreateOrderLogic) TryCreateOrder(in *order.TryCreateOrderReq) (*order.TryCreateOrderResp, error) {
	// todo: add your logic here and delete this line
	var shippingId int64

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

	//

	// ===== 第三步：在屏障中执行 Try 逻辑 =====

	// 开启子事务屏障（屏障会保证在数据库事务中执行，并且自动处理重复请求的幂等性）
	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {

		// 1. 检查订单是否已存在（幂等性处理）
		existing, _ := l.svcCtx.OrdersModel.FindOne(l.ctx, in.OrderId)
		if existing != nil {
			// 已存在则直接返回成功（幂等）
			return nil
		}

		//插入信息
		//ship
		var insertShipping model.Shipping
		insertShipping.Orderid = in.OrderId
		insertShipping.Userid = uint64(in.UserId)
		insertShipping.ReceiverName = in.ReceiverName
		insertShipping.ReceiverPhone = in.ReceiverPhone
		insertShipping.ReceiverMobile = in.ReceiverPhone
		insertShipping.ReceiverProvince = in.ReceiverProvince
		insertShipping.ReceiverCity = in.ReceiverCity
		insertShipping.ReceiverDistrict = in.ReceiverDistrict
		insertShipping.ReceiverAddress = in.ReceiverDetailAddress

		insertShippingRes, err := l.svcCtx.ShippingModel.TxInsert(tx, &insertShipping)
		if err != nil {
			return err
		}
		shippingId, err = insertShippingRes.LastInsertId()
		if err != nil {
			return err
		}
		//orderitem

		insertOrderitem := &model.Orderitem{

			OrderId:      in.OrderId,           // 订单id
			UserId:       uint64(in.UserId),    // 用户id
			ProductId:    uint64(in.ProductId), // 商品id
			ProductName:  in.ProductName,       // 商品名称
			ProductImage: in.ProductImage,      // 商品图片地址
			CurrentPrice: in.CurrentPrice,      // 生成订单时的商品单价，单位是元,保留两位小数
			Quantity:     in.Quantity,          // 商品数量

		}

		insertOrderitemRes, err := l.svcCtx.OrderitemModel.TxInsert(tx, insertOrderitem)
		if err != nil {
			return err
		}

		_, err = insertOrderitemRes.LastInsertId()
		if err != nil {
			return err
		}
		//order
		insertOrder := model.Orders{
			Id:         in.OrderId,
			Userid:     uint64(in.UserId),
			Shoppingid: shippingId,
			Postage:    in.Postage,
			Status:     0, // 0=待确认
		}

		_, err = l.svcCtx.OrdersModel.TxInsert(tx, &insertOrder)
		if err != nil {
			return err // 插入失败则回滚事务
		}

		return nil

	}); err != nil {
		// 如果事务执行失败，包装错误并返回
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "create new order Exception : err: %v", err)
	}

	return &order.TryCreateOrderResp{
		ShippingId: shippingId,
	}, nil
}

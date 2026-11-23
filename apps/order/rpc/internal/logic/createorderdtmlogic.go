package logic

import (
	"context"

	"lebron/apps/order/rpc/internal/svc"
	"lebron/apps/order/rpc/order"
	"lebron/apps/product/rpc/productclient"
	"lebron/apps/user/rpc/userclient"
	"lebron/pkg/snowflake"
	"lebron/pkg/xerr"

	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type CreateOrderDTMLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderDTMLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderDTMLogic {
	return &CreateOrderDTMLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DTM（分布式事务管理器）正向操作：创建订单（包含完整订单信息）
func (l *CreateOrderDTMLogic) CreateOrderDTM(in *order.AddOrderReq) (*order.AddOrderResp, error) {
	// todo: add your logic here and delete this line

	var (
		userRpcRes        *userclient.GetUserResponse
		receiveAddressRes *userclient.GetAddressResponse
		productRpcRes     *productclient.ProductItem
	)

	//并发校验校验产品 收货信息 用户是否存在

	checkUser := func() error {
		var err error
		var userreq userclient.GetUserRequest
		userreq.Id = in.Userid
		userRpcRes, err = l.svcCtx.UserRpc.GetUser(l.ctx, &userreq)
		if err != nil {
			return nil
		}
		return nil
	}
	checkproduct := func() error {
		var err error
		var productreq productclient.ProductReq
		productreq.ProductId = in.Productid
		productRpcRes, err = l.svcCtx.ProductRpc.GetProductDetail(l.ctx, &productreq)
		if err != nil {
			return nil
		}
		return nil
	}
	checkreceiveAddress := func() error {
		var err error
		var userAddressreq userclient.GetAddressRequest
		userAddressreq.Id = in.ReceiveAddressId
		receiveAddressRes, err = l.svcCtx.UserRpc.GetAddress(l.ctx, &userAddressreq)
		if err != nil {
			return nil
		}
		return nil
	}
	// 并发执行上述三个校验任务（使用 go-zero 的 mr.Finish 工具）

	_ = mr.Finish(checkproduct, checkUser, checkreceiveAddress)
	// ---------- 校验结果处理 ----------

	if userRpcRes == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DataNoExistError), "error! user not exist exception : %+v  ", userRpcRes)
	}
	if receiveAddressRes == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DataNoExistError), "error! receiveAddressRes not exist exception : %+v  ", receiveAddressRes)
	}

	if productRpcRes == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DataNoExistError), "error! productRpcRes not exist exception : %+v  ", productRpcRes)
	}

	if productRpcRes.Stock <= 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("product understock"), "product understock")
	}

	// ---------- 生成订单ID ----------
	// 使用雪花算法生成全局唯一订单ID（字符串形式）
	orderId := snowflake.GenIDString()

	// ---------- 分布式事务准备 ----------
	// //databasde
	// db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	// if err != nil {
	// 	return nil, status.Error(501, err.Error())
	// }

	// ===== 第三步：获取 DTM 服务地址 =====
	dtmServer := l.svcCtx.Config.DTM.Server

	// ===== 第四步：获取 Order 服务自身的 gRPC 地址（供 DTM 回调）=====
	// orderGrpcServer, _ := l.svcCtx.Config.OrderRpc.BuildTarget()

	orderGrpcServer := l.svcCtx.Config.ListenOn // 注意：需确保外部可访问，不要用 0.0.0.0
	if orderGrpcServer == "" {
		orderGrpcServer = "127.0.0.1:8080" // 可从环境变量读取
	}

	productGrpcServer := l.svcCtx.Config.ProductRpc.Target

	// ===== 第五步：获取 Product 服务的 gRPC 地址 =====

	// ===== 第六步：构造 Try 请求参数 =====

	// Order 服务的 Try 请求
	orderTryReq := &order.TryCreateOrderReq{
		OrderId:               orderId,
		UserId:                in.Userid,
		ProductId:             in.Productid,
		Quantity:              in.Quantity,
		Postage:               in.Postage,
		ReceiveAddressId:      in.ReceiveAddressId,
		ReceiverName:          receiveAddressRes.Address.Name,
		ReceiverPhone:         receiveAddressRes.Address.Phone,
		ReceiverProvince:      receiveAddressRes.Address.Province,
		ReceiverCity:          receiveAddressRes.Address.City,
		ReceiverDistrict:      receiveAddressRes.Address.Region,
		ReceiverDetailAddress: receiveAddressRes.Address.DetailAddress,
		ProductName:           productRpcRes.Name,
		ProductImage:          productRpcRes.ImageUrl,
		CurrentPrice:          productRpcRes.Price,
	}

	// Product 服务的 Try 请求
	productTryReq := &productclient.TryDeductStockReq{
		ProductId: in.Productid,
		Quantity:  in.Quantity,
	}

	// ===== 第七步：发起 TCC 全局事务 =====
	gid := dtmgrpc.MustGenGid(dtmServer)

	err := dtmgrpc.TccGlobalTransaction(dtmServer, gid, func(tcc *dtmgrpc.TccGrpc) error {

		// ---- 分支1：Order 服务（创建订单数据）----
		targetorder := orderGrpcServer // 例如 "localhost:8080"
		// CallBranch 参数：请求数据，Try路径，Confirm路径，Cancel路径
		// 构造 gRPC 方法路径
		tryURL := targetorder + "/order.Order/TryCreateOrder"
		confirmURL := targetorder + "/order.Order/ConfirmCreateOrder"
		cancelURL := targetorder + "/order.Order/CancelCreateOrder"

		// 注册分支
		err := tcc.CallBranch(orderTryReq, tryURL, confirmURL, cancelURL, &order.TryCreateOrderResp{})
		if err != nil {
			return err
		}

		targetproduct := productGrpcServer // 例如 "localhost:8080"
		// CallBranch 参数：请求数据，Try路径，Confirm路径，Cancel路径
		// 构造 gRPC 方法路径
		tryproductURL := targetproduct + "/product.Product/TryDeductStock"
		confirmproductURL := targetproduct + "/product.Product/ConfirmDeductStock"
		cancelproductURL := targetproduct + "/product.Product/CancelDeductStock"

		// 注册分支
		err = tcc.CallBranch(
			productTryReq, tryproductURL, confirmproductURL, cancelproductURL, &productclient.TryDeductStockResp{})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		// 如果事务执行失败，包装错误并返回
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "create new order Exception : err: %v", err)
	}

	//---------- 返回成功响应 ----------

	return &order.AddOrderResp{
		Id: orderId, // 返回生成的订单ID
	}, nil
}

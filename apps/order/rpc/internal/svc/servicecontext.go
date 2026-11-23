package svc

import (
	"lebron/apps/order/rpc/internal/config"
	"lebron/apps/order/rpc/internal/model"
	"lebron/apps/order/rpc/orderclient"
	"lebron/apps/product/rpc/productclient"
	"lebron/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	OrdersModel    model.OrdersModel
	OrderitemModel model.OrderitemModel
	ShippingModel  model.ShippingModel
	ProductRpc     productclient.Product
	UserRpc        userclient.User
	OrderRpc       orderclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:         c,
		OrdersModel:    model.NewOrdersModel(sqlConn, c.Cache),
		OrderitemModel: model.NewOrderitemModel(sqlConn, c.Cache),
		ShippingModel:  model.NewShippingModel(sqlConn, c.Cache),
		UserRpc:        userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ProductRpc:     productclient.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
	}
}

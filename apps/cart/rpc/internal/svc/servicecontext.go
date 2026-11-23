package svc

import (
	"lebron/apps/cart/rpc/internal/config"
	"lebron/apps/cart/rpc/internal/model"
	"lebron/apps/product/rpc/productclient"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	CartModel  model.CartModel
	ProductRpc productclient.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:    c,
		CartModel: model.NewCartModel(sqlConn, c.Cache),

		ProductRpc: productclient.NewProduct(zrpc.MustNewClient(c.ProductRpcConf)),
	}
}

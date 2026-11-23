package svc

import (
	"lebron/apps/reply/rpc/internal/config"
	"lebron/apps/reply/rpc/internal/model"
	"lebron/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	ReplyModel  model.ReplyModel
	UserRpcConf userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:      c,
		ReplyModel:  model.NewReplyModel(sqlConn, c.Cache),
		UserRpcConf: userclient.NewUser(zrpc.MustNewClient(c.UserRpcConf)),
	}
}

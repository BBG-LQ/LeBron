package svc

import (
	"lebron/apps/user/rpc/internal/config"
	"lebron/apps/user/rpc/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                  config.Config
	UserModel               model.UserModel
	UserCollectionModel     model.UserCollectionModel
	UserReceiveAddressModel model.UserReceiveAddressModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:                  c,
		UserModel:               model.NewUserModel(sqlConn, c.Cache),
		UserCollectionModel:     model.NewUserCollectionModel(sqlConn, c.Cache),
		UserReceiveAddressModel: model.NewUserReceiveAddressModel(sqlConn, c.Cache),
	}
}

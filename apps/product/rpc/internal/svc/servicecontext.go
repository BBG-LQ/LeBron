package svc

import (
	"lebron/apps/product/rpc/internal/config"
	"lebron/apps/product/rpc/model"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"golang.org/x/sync/singleflight"
)

const localCacheExpire = time.Duration(time.Second * 60)

type ServiceContext struct {
	Config                config.Config
	ProductModel          model.ProductModel
	CategoryModel         model.CategoryModel
	ProductOperationModel model.ProductOperationModel
	BizRedis              *redis.Redis
	SingleGroup           singleflight.Group
	LocalCache            *collection.Cache
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	LocalCache, _ := collection.NewCache(localCacheExpire)
	return &ServiceContext{
		Config:                c,
		ProductModel:          model.NewProductModel(sqlConn, c.Cache),
		CategoryModel:         model.NewCategoryModel(sqlConn, c.Cache),
		ProductOperationModel: model.NewProductOperationModel(sqlConn, c.Cache),
		BizRedis:              redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		SingleGroup:           singleflight.Group{},
		LocalCache:            LocalCache,
	}
}

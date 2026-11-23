package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		DataSource string
	}
	DTM struct {
		Server string
	}

	Cache      cache.CacheConf
	UserRpc    zrpc.RpcClientConf
	ProductRpc zrpc.RpcClientConf
	OrderRpc   zrpc.RpcClientConf
}

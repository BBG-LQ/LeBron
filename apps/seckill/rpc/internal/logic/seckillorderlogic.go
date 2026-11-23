package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"lebron/apps/product/rpc/productclient"
	"lebron/apps/seckill/rpc/internal/svc"
	"lebron/apps/seckill/rpc/seckill"
	"lebron/pkg/batcher"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SeckillOrderLogic 秒杀订单逻辑结构体
// 包含上下文、服务上下文、日志、限流器、本地缓存、批处理器
type SeckillOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger

	limiter    *limit.PeriodLimit
	localCache *collection.Cache
	batcher    *batcher.Batcher
}

// KafkaData 定义要发送到Kafka的数据结构，包含用户ID和商品ID
type KafkaData struct {
	Uid int64 `json:"uid"`
	Pid int64 `json:"pid"`
}

// 常量定义：限流周期（秒）、限流配额（每次周期允许请求数）、用户限流key前缀、本地缓存过期时间、
// 批处理器参数（大小、缓冲、工作协程数、间隔）
const (
	limitPeriod       = 10
	limitQuota        = 1
	seckillUserPrefix = "seckill#u#"
	localCacheExpire  = time.Second * 60

	batcherSize     = 100
	batcherBuffer   = 100
	batcherWorker   = 10
	batcherInterval = time.Second
)

func NewSeckillOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillOrderLogic {

	// 创建本地缓存，过期时间为60秒
	localCache, err := collection.NewCache(localCacheExpire)
	if err != nil {
		panic(err)
	}

	s := &SeckillOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),

		localCache: localCache,
		// 初始化限流器：周期10秒，配额1次，使用Redis，key前缀为seckill#u#
		limiter: limit.NewPeriodLimit(limitPeriod, limitQuota, &svcCtx.BizRedis, seckillUserPrefix),
	}

	// 创建批处理器，配置大小、缓冲、工作协程数、间隔
	b := batcher.New(
		batcher.WithSize(batcherSize),
		batcher.WithBuffer(batcherBuffer),
		batcher.WithWorker(batcherWorker),
		batcher.WithInterval(batcherInterval),
	)

	// 设置分片函数：根据商品ID字符串转换为int64，然后对工作协程数取模，决定由哪个worker处理
	b.Sharding = func(key string) int {
		pid, _ := strconv.ParseInt(key, 10, 64)
		return int(pid) % batcherWorker
	}

	// 设置批处理执行函数：将收集到的数据批量序列化为JSON并推送到Kafka
	b.Do = func(ctx context.Context, val map[string][]interface{}) {
		var msgs []*KafkaData
		// 遍历所有分片的值（每个分片是一个[]interface{}）
		for _, vs := range val {
			for _, v := range vs {
				// 断言为*KafkaData并追加到msgs切片
				msgs = append(msgs, v.(*KafkaData))
			}
		}
		// 将msgs序列化为JSON
		kd, err := json.Marshal(msgs)
		if err != nil {
			logx.Errorf("Batcher.Do json.Marshal msgs: %v error: %v", msgs, err)
		}
		// 将JSON字符串推送到Kafka
		if err = s.svcCtx.KafkaPusher.Push(ctx, string(kd)); err != nil {
			logx.Errorf("KafkaPusher.Push kd: %s error: %v", string(kd), err)
		}
	}

	s.batcher = b
	s.batcher.Start() // 启动批处理器的后台协程

	return s
}

// 秒杀订单
// SeckillOrder 处理秒杀订单请求
// 1. 对用户进行限流检查
// 2. 调用商品RPC检查库存
// 3. 将订单数据加入批处理器，异步发送到Kafka
// 返回空响应或错误
func (l *SeckillOrderLogic) SeckillOrder(in *seckill.SeckillOrderRequest) (*seckill.SeckillOrderResponse, error) {
	// todo: add your logic here and delete this line

	// 对用户ID进行限流，Take返回状态码和剩余时间（忽略剩余时间）
	code, _ := l.limiter.Take(strconv.FormatInt(in.UserId, 10))
	if code == limit.OverQuota {
		return nil, status.Errorf(codes.OutOfRange, "Number of requests exceeded the limit")
	}

	// 调用商品RPC服务获取商品信息
	pro, err := l.svcCtx.ProductRPC.GetProductDetail(l.ctx, &productclient.ProductReq{ProductId: in.ProductId})
	if err != nil {
		return nil, err
	}
	//库存是否充足

	if pro.Stock <= 0 {
		return nil, status.Errorf(codes.OutOfRange, "Insufficient stock")
	}

	// 将订单数据（用户ID、商品ID）加入批处理器，key为商品ID字符串，用于分片

	if err = l.batcher.Add(strconv.FormatInt(in.ProductId, 10), KafkaData{Uid: in.UserId, Pid: in.ProductId}); err != nil {
		// 加入失败则记录错误日志（但不影响返回成功，因为异步处理，这里可能根据业务决定是否返回错误）
		logx.Errorf("l.batcher.Add uid: %d pid: %d error: %v", in.UserId, in.ProductId, err)
	}

	return &seckill.SeckillOrderResponse{}, nil
}

package logic

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"lebron/apps/order/rpc/internal/svc"
	"lebron/apps/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建订单（简单版本，仅包含用户ID和商品ID）
func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	// todo: add your logic here and delete this line

	oid := genOrderID(time.Now())

	err := l.svcCtx.OrdersModel.CreateOrder(l.ctx, oid, in.Uid, in.Pid)
	if err != nil {
		logx.Errorf("OrderModel.CreateOrder oid: %s uid: %d pid: %d", oid, in.Uid, in.Pid)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &order.CreateOrderResponse{}, nil
}

var num int64

func genOrderID(t time.Time) string {
	//时间t
	s := t.Format("20200615121212")

	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3

	ms := sup(m, 3)
	// 进程pid
	p := os.Getpid() % 1000
	ps := sup(int64(p), 3)
	//原子性地将全局计数器 num 加1，并返回新值（旧值+1）
	i := atomic.AddInt64(&num, 1)
	r := i % 10000
	rs := sup(r, 4)
	// 拼接：时间（秒） + 毫秒（3位） + 进程ID（3位） + 序号（4位），共 14+3+3+4 = 24 位
	n := fmt.Sprintf("%s%s%s%s", s, ms, ps, rs)
	return n

}

func sup(i int64, n int) string {

	m := fmt.Sprintf("%d", i)

	if len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}

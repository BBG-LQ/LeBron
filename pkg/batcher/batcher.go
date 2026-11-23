package batcher

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

// ErrFull 表示通道已满，无法添加新消息
var ErrFull = errors.New("channel is full")

// Option 定义配置选项的接口，所有选项都必须实现 apply 方法
type Option interface {
	apply(*options)
}

// options 存放 Batcher 的内部配置参数
type options struct {
	size     int           // 每次触发 Do 时累积的消息数量阈值
	buffer   int           // 每个 worker 通道的缓冲大小
	worker   int           // 并发 worker 的数量（也即分片数）
	interval time.Duration // 定时触发 Do 的时间间隔
}

// check 对配置进行校验，若未设置则赋予默认值
func (o options) check() {
	if o.size <= 0 {
		o.size = 100
	}
	if o.buffer <= 0 {
		o.buffer = 100
	}
	if o.worker <= 0 {
		o.worker = 5
	}
	if o.interval <= 0 {
		o.interval = time.Second
	}
}

// funcOption 是一个函数类型的 Option，通过闭包实现 apply
type funcOption struct {
	f func(*options)
}

// apply 执行存储的函数来修改 options
func (fo *funcOption) apply(o *options) {
	fo.f(o)
}

// newOption 构造一个 funcOption
func newOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithSize 设置每次批量处理的消息数量阈值
func WithSize(s int) Option {
	return newOption(func(o *options) {
		o.size = s
	})
}

// WithBuffer 设置每个 worker 通道的缓冲区大小
func WithBuffer(b int) Option {
	return newOption(func(o *options) {
		o.buffer = b
	})
}

// WithWorker 设置 worker 的数量（并发数/分片数）
func WithWorker(w int) Option {
	return newOption(func(o *options) {
		o.worker = w
	})
}

// WithInterval 设置定时刷新的时间间隔
func WithInterval(i time.Duration) Option {
	return newOption(func(o *options) {
		o.interval = i
	})
}

// msg 表示一条待处理的消息，包含分组 key 和值
type msg struct {
	key string
	val interface{}
}

// Batcher 是批量处理器的核心结构
type Batcher struct {
	opts options

	// Do 是用户提供的批量处理函数，接收上下文和按 key 分组的 value 切片
	Do func(ctx context.Context, val map[string][]interface{})
	// Sharding 是用户提供的分片函数，根据 key 返回一个整数，用于选择 worker
	Sharding func(key string) int
	chans    []chan *msg    // 每个 worker 对应的消息通道
	wait     sync.WaitGroup // 用于等待所有 worker 协程退出
}

// New 创建一个 Batcher 实例，应用传入的 Option 配置
func New(opts ...Option) *Batcher {
	b := &Batcher{}
	for _, opt := range opts {
		opt.apply(&b.opts)
	}
	b.opts.check() // 设置默认值

	b.chans = make([]chan *msg, b.opts.worker)
	for i := 0; i < b.opts.worker; i++ {
		b.chans[i] = make(chan *msg, b.opts.buffer)
	}
	return b
}

// Start 启动所有 worker 协程，必须在使用 Add 之前调用
func (b *Batcher) Start() {
	if b.Do == nil {
		log.Fatal("Batcher: Do func is nil")
	}
	if b.Sharding == nil {
		log.Fatal("Batcher: Sharding func is nil")
	}
	b.wait.Add(len(b.chans))
	for i, ch := range b.chans {
		go b.merge(i, ch) // 每个 worker 启动一个 merge 协程
	}
}

// Add 向 Batcher 添加一条消息，根据 key 分片后发送到对应的通道
// 若通道已满则返回 ErrFull，否则返回 nil
func (b *Batcher) Add(key string, val interface{}) error {
	ch, msg := b.add(key, val)
	select {
	case ch <- msg:
	default:
		return ErrFull
	}
	return nil
}

// add 根据 key 计算目标 worker 索引并构造 msg，返回目标通道和消息对象
func (b *Batcher) add(key string, val interface{}) (chan *msg, *msg) {
	sharding := b.Sharding(key) % b.opts.worker
	ch := b.chans[sharding]
	msg := &msg{key: key, val: val}
	return ch, msg
}

// merge 是每个 worker 的主循环，从通道读取消息，累积到一定数量或超时后调用 Do 处理
// idx 是 worker 索引，用于错开各 worker 的定时器相位，避免同时触发
func (b *Batcher) merge(idx int, ch <-chan *msg) {
	defer b.wait.Done()

	var (
		msg        *msg
		count      int    // 当前累积的消息总数（用于判断是否达到 size 阈值）
		closed     bool   // 标记通道是否已收到关闭信号（nil 消息）
		lastTicker = true // 标识是否为第一次启动的 ticker（用于重新设置 interval）
		interval   = b.opts.interval
		vals       = make(map[string][]interface{}, b.opts.size) // 按 key 分组存储 value
	)
	// 为每个 worker 设置不同的初始间隔，使它们的首次触发时间错开，减少并发压力
	if idx > 0 {
		interval = time.Duration(int64(idx) * (int64(b.opts.interval) / int64(b.opts.worker)))
	}
	ticker := time.NewTicker(interval)
	for {
		select {
		case msg = <-ch:
			if msg == nil { // nil 作为关闭信号
				closed = true
				break
			}
			count++
			vals[msg.key] = append(vals[msg.key], msg.val)
			if count >= b.opts.size {
				break // 达到阈值后退出 select，执行下方的处理
			}
			continue // 未达阈值则继续等待下一条消息
		case <-ticker.C:
			if lastTicker {
				// 第一次 ticker 触发后，重置为标准间隔，并标记为非首次
				ticker.Stop()
				ticker = time.NewTicker(b.opts.interval)
				lastTicker = false
			}
		}
		// 当有数据时执行 Do 处理
		if len(vals) > 0 {
			ctx := context.Background()
			b.Do(ctx, vals)
			// 重置 map 和计数器
			vals = make(map[string][]interface{}, b.opts.size)
			count = 0
		}
		if closed {
			ticker.Stop()
			return
		}
	}
}

// Close 关闭 Batcher，向所有通道发送 nil 信号通知 worker 退出，并等待所有 worker 完成
func (b *Batcher) Close() {
	for _, ch := range b.chans {
		ch <- nil // 发送关闭信号
	}
	b.wait.Wait() // 等待所有 merge 协程结束
}

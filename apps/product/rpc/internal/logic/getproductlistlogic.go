package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"lebron/apps/product/rpc/internal/svc"
	"lebron/apps/product/rpc/model"
	"lebron/apps/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetProductlistLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductlistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductlistLogic {
	return &GetProductlistLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 缓存常量
const (
	defaultPageSize = 10
	defaultLimit    = 300
	expireTime      = 3600 * 24 * 3
)

// ProductList 查询商品列表
func (l *GetProductlistLogic) GetProductlist(in *product.ProductListRequest) (*product.ProductListResponse, error) {
	// todo: add your logic here and delete this line

	// 判断是否存在此分类

	_, err := l.svcCtx.CategoryModel.FindOne(l.ctx, uint64(in.CategoryId))
	if err == model.ErrNotFound {
		return nil, status.Error(codes.NotFound, "分类不在")
	}

	if in.Cursor == 0 {
		in.Cursor = time.Now().Unix()
	}
	if in.Ps == 0 {
		in.Ps = 2
	}

	var (
		isCache, isEnd   bool
		lastID, lastTime int64
		firstPage        []*product.ProductItem
		products         []*model.Product
	)

	//尝试从redis有序集合中读取该分类下的商品ID列表（按时间倒序，游标为时间戳）
	pids, err := l.cacheProductList(l.ctx, in.CategoryId, in.Cursor, int64(in.Ps))

	if len(pids) == int(in.Ps) {

		isCache = true
		if pids[len(pids)-1] == -1 {
			isEnd = true
		}

		products, err = l.productsByIds(l.ctx, pids)

		if err != nil {
			return nil, err
		}

		for _, p := range products {
			firstPage = append(firstPage, &product.ProductItem{
				ProductId:   int64(p.Id),
				Name:        p.Name,
				CreateTime:  p.CreateTime.Unix(),
				Description: "1",
			})

		}

	} else {

		var (
			err   error
			ctime = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")
		)
		// 从数据库中查询该分类下的商品列表（按时间倒序，游标为时间戳）
		products, err = l.svcCtx.ProductModel.CategoryProducts(l.ctx, ctime, int64(in.CategoryId), defaultLimit)

		if err != nil {
			return nil, err
		}

		//// 确定第一页展示的数据：若查询结果多于请求页大小，则取前 in.Ps 条，

		var firstPageProducts []*model.Product
		if len(products) > int(in.Ps) {
			firstPageProducts = products[:int(in.Ps)]

		} else {
			firstPageProducts = products
			isEnd = true
		}

		//转换为响应格式
		for _, p := range firstPageProducts {

			firstPage = append(firstPage, &product.ProductItem{
				ProductId:   int64(p.Id),
				Name:        p.Name,
				CreateTime:  p.CreateTime.Unix(),
				Description: "2",
			})
		}

	}

	//

	//如果 firstPage 非空，计算最后一条记录的 ID 和时间戳（用于返回给客户端作为下一页游标）

	if len(firstPage) > 0 {
		pageLast := firstPage[len(firstPage)-1]
		lastID = pageLast.ProductId
		lastTime = pageLast.CreateTime

		if lastTime < 0 {
			lastTime = 0
		}

		// 去重处理：如果客户端传入了 productId（用于断点续传），则跳过之前已返回的记录
		// 这里假设前端会传入上次返回的最后一条记录的 productId 和 cursor，用于避免重复

		for k, p := range firstPage {
			if p.CreateTime == in.Cursor && p.ProductId == in.ProductId {
				firstPage = firstPage[k:] // 截取从匹配项开始
				break
			}
		}

	}

	ret := &product.ProductListResponse{
		IsEnd:     isEnd,
		Timestamp: lastTime,
		ProductId: lastID,
		Products:  firstPage,
	}

	if !isCache {
		threading.GoSafe(func() {
			// 如果查询结果不足 defaultLimit 且大于0，说明没有更多数据，添加一个结束标记（ID=-1，时间为零值）
			if len(products) < defaultLimit && len(products) > 0 {

				endTime, _ := time.Parse("2006-01-02 15:04:05", "0-0-0 00:00:00")

				products = append(products, &model.Product{
					Id:         999,
					CreateTime: endTime,
				})

			}

			_ = l.addCacheProductList(context.Background(), products)
		})
	}

	return ret, nil
}

// cacheProductList 从 Redis 有序集合中按分数（时间戳）倒序获取指定分类的商品ID列表
func (l *GetProductlistLogic) cacheProductList(ctx context.Context, cid int32, cursor, ps int64) ([]int64, error) {

	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, categoryKey(cid), 0, cursor, 0, int(ps))

	// fmt.Println("*******************")
	// fmt.Println("数据categoryKey(cid)：%v", categoryKey(cid))
	// fmt.Println("数据cursor%v：", cursor)
	// fmt.Println("数据%v：", pairs)

	// fmt.Println("数据int(ps：%v：", int(ps))
	// fmt.Println("*******************")
	if err != nil {
		return nil, err
	}

	var ids []int64

	for _, pair := range pairs {
		id, _ := strconv.ParseInt(pair.Key, 10, 64)
		ids = append(ids, id)

	}

	return ids, nil

}

// categoryKey 生成分类缓存的 Redis 键
func categoryKey(cid int32) string {
	return fmt.Sprintf("category:%d", cid)
}

// productsByIds 根据商品ID列表并发查询商品详情，返回商品切片
func (l *GetProductlistLogic) productsByIds(ctx context.Context, pids []int64) ([]*model.Product, error) {

	products, err := mr.MapReduce[int64, *model.Product, []*model.Product](

		func(source chan<- int64) {
			for _, pid := range pids {
				source <- pid
			}
		}, func(item int64, writer mr.Writer[*model.Product], cancel func(error)) {

			p, err := l.svcCtx.ProductModel.FindOne(ctx, uint64(item))

			if err != nil {
				cancel(err)
				return
			}

			writer.Write(p)

		}, func(pipe <-chan *model.Product, writer mr.Writer[[]*model.Product], cancel func(error)) {

			var ps []*model.Product

			for item := range pipe {

				ps = append(ps, item)
			}

			writer.Write(ps)
		},
	)

	if err != nil {
		return nil, err
	}

	return products, nil

}

// addCacheProductList 将商品列表写入 Redis 有序集合，用于缓存商品列表
// 每个商品以 ID 为成员，创建时间戳为分数，并设置过期时间
func (l *GetProductlistLogic) addCacheProductList(ctx context.Context, products []*model.Product) error {

	if len(products) == 0 {
		return nil
	}

	for _, p := range products {

		score := p.CreateTime.Unix()
		if score < 0 {
			score = 0
		}
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, categoryKey(int32(p.Cateid)), score, strconv.Itoa(int(p.Id)))
		if err != nil {
			return err
		}
	}

	return l.svcCtx.BizRedis.ExpireCtx(ctx, categoryKey(int32(products[0].Cateid)), expireTime)

}

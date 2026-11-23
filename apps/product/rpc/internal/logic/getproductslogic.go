package logic

import (
	"context"
	"strconv"
	"strings"

	"lebron/apps/product/rpc/internal/svc"
	"lebron/apps/product/rpc/model"
	"lebron/apps/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type GetProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductsLogic {
	return &GetProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Products 批量获取商品信息。
func (l *GetProductsLogic) GetProducts(in *product.ProductsRequest) (*product.ProductsResponse, error) {
	// todo: add your logic here and delete this line

	products := make(map[int64]*product.ProductItem)
	pids := strings.Split(in.ProductIds, ",")

	ps, err := mr.MapReduce[string, *model.Product, []*model.Product](
		func(source chan<- string) {

			for _, pid := range pids {
				source <- pid
			}

		},
		func(item string, writer mr.Writer[*model.Product], cancel func(error)) {

			pid, err := strconv.ParseInt(item, 10, 64)

			if err != nil {

				return
			}
			p, err := l.svcCtx.ProductModel.FindOne(l.ctx, uint64(pid))

			if err != nil {
				return
			}

			writer.Write(p)

		},
		func(pipe <-chan *model.Product, writer mr.Writer[[]*model.Product], cancel func(error)) {

			var r []*model.Product

			for p := range pipe {
				r = append(r, p)
			}

			writer.Write(r)

		})

	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		products[int64(p.Id)] = &product.ProductItem{
			ProductId: int64(p.Id),
			Name:      p.Name,
			Price:     p.Price,
		}

	}

	return &product.ProductsResponse{
		Products: products,
	}, nil

}

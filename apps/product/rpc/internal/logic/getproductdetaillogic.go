package logic

import (
	"context"
	"fmt"

	"lebron/apps/product/rpc/internal/svc"
	"lebron/apps/product/rpc/model"
	"lebron/apps/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductDetailLogic {
	return &GetProductDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Product 获取单个商品详情
func (l *GetProductDetailLogic) GetProductDetail(in *product.ProductReq) (*product.ProductItem, error) {
	// todo: add your logic here and delete this line

	v, err, _ := l.svcCtx.SingleGroup.Do(fmt.Sprintf("product:%d", in.ProductId), func() (interface{}, error) {
		return l.svcCtx.ProductModel.FindOne(l.ctx, uint64(in.GetProductId()))
	})

	if err != nil {
		return nil, err
	}

	p := v.(*model.Product)

	return &product.ProductItem{
		ProductId:   int64(p.Id),
		Name:        p.Name,
		Description: p.Detail,
		ImageUrl:    p.Images,
		Stock:       p.Stock,
		Cateid:      int64(p.Cateid),
		Price:       p.Price,
		Status:      p.Status,
	}, nil
}

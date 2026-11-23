package logic

import (
	"context"
	"fmt"

	"lebron/apps/product/rpc/internal/svc"
	"lebron/apps/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperationProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	ProductlistLogic *GetProductlistLogic
}

func NewOperationProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationProductsLogic {
	return &OperationProductsLogic{
		ctx:              ctx,
		svcCtx:           svcCtx,
		Logger:           logx.WithContext(ctx),
		ProductlistLogic: NewGetProductlistLogic(ctx, svcCtx),
	}
}

const (
	validStatus          = 1
	operationProductsKey = "operation#products"
)

func (l *OperationProductsLogic) OperationProducts(in *product.OperationProductsRequest) (*product.OperationProductsResponse, error) {
	// todo: add your logic here and delete this line

	//本地缓存
	opProducts, ok := l.svcCtx.LocalCache.Get(operationProductsKey)

	fmt.Println("*******************", ok)
	if ok {

		fmt.Println("*******************")
		fmt.Println("*******************")
		return &product.OperationProductsResponse{
			Products: opProducts.([]*product.ProductItem),
		}, nil
	}

	//数据库

	pos, err := l.svcCtx.ProductOperationModel.OperationProducts(l.ctx, validStatus)

	if err != nil {
		return nil, err
	}

	var pids []int64
	for _, p := range pos {
		pids = append(pids, int64(p.ProductId))

	}

	products, err := l.ProductlistLogic.productsByIds(l.ctx, pids)

	if err != nil {
		return nil, err
	}

	var pItem []*product.ProductItem

	for _, p := range products {

		pItem = append(pItem, &product.ProductItem{
			ProductId: int64(p.Id),
			Name:      p.Name,
		})

	}

	// 存入本地缓存

	l.svcCtx.LocalCache.Set(operationProductsKey, pItem)

	//

	return &product.OperationProductsResponse{
		Products: pItem,
	}, nil
}

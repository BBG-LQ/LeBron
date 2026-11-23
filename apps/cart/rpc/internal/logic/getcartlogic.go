package logic

import (
	"context"
	"strconv"
	"strings"

	"lebron/apps/cart/rpc/cart"
	"lebron/apps/cart/rpc/internal/model"
	"lebron/apps/cart/rpc/internal/svc"
	"lebron/apps/product/rpc/product"
	"lebron/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartLogic {
	return &GetCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户的完整购物车信息（包含商品详情）
func (l *GetCartLogic) GetCart(in *cart.GetCartReq) (*cart.GetCartResp, error) {
	// todo: add your logic here and delete this line

	listCarts, err := l.svcCtx.CartModel.FindAllByUid(l.ctx, uint64(in.UserId))

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "Failed  get user's cart list err : %v , in :%+v", err, in)
	}

	var listPro []int64

	qui := make(map[int64]int64)

	chect := make(map[int64]bool)

	for _, v := range listCarts {
		listPro = append(listPro, int64(v.Productid))
		qui[int64(v.Productid)] = v.Quantity

		if v.Checked == 0 {
			chect[int64(v.Productid)] = false
		} else {
			chect[int64(v.Productid)] = true
		}

	}
	//
	var builder strings.Builder
	for i, id := range listPro {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(strconv.FormatInt(id, 10))
	}
	productIdsStr := builder.String()

	getprodsRes := &product.ProductsRequest{
		ProductIds: productIdsStr,
	}

	pros, err := l.svcCtx.ProductRpc.GetProducts(l.ctx, getprodsRes)

	var produs []*cart.CartItem
	for _, v := range pros.Products {

		item := &cart.CartItem{
			ProductId:   v.ProductId,
			ProductName: v.Name,
			Price:       v.Price,
			Quantity:    qui[v.ProductId],
			Checked:     chect[v.ProductId],
		}

		produs = append(produs, item)

	}

	return &cart.GetCartResp{

		Cart: &cart.Cart{
			UserId: in.UserId,
			Items:  produs,
		},
	}, nil
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"lebron/apps/app/api/internal/svc"
	"lebron/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlashSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 限时抢购
func NewFlashSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlashSaleLogic {
	return &FlashSaleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlashSaleLogic) FlashSale() (resp *types.FlashSaleResponse, err error) {
	// todo: add your logic here and delete this line

	return
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"lebron/apps/app/api/internal/svc"
	"lebron/apps/app/api/internal/types"
	"lebron/apps/user/rpc/user"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type EditReceiveAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// edit user receiveAddress
func NewEditReceiveAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditReceiveAddressLogic {
	return &EditReceiveAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditReceiveAddressLogic) EditReceiveAddress(req *types.UserReceiveAddressEditReq) (resp *types.UserReceiveAddressEditRes, err error) {
	// todo: add your logic here and delete this line
	var editAddreq user.EditAddressRequest

	coprErr := copier.Copy(&editAddreq, req)
	if coprErr != nil {
		return nil, coprErr
	}

	_, err = l.svcCtx.UserRPC.EditAddress(l.ctx, &editAddreq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

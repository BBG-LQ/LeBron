// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"lebron/apps/app/api/internal/svc"
	"lebron/apps/app/api/internal/types"
	"lebron/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelReceiveAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// del user receiveAddress list
func NewDelReceiveAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelReceiveAddressLogic {
	return &DelReceiveAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelReceiveAddressLogic) DelReceiveAddress(req *types.UserReceiveAddressDelReq) (resp *types.UserReceiveAddressDelRes, err error) {
	// todo: add your logic here and delete this line

	var delAddreq user.DeleteAddressRequest
	delAddreq.Id = req.Id

	_, err = l.svcCtx.UserRPC.DeleteAddress(l.ctx, &delAddreq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

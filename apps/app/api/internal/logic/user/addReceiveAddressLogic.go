// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"encoding/json"

	"lebron/apps/app/api/internal/svc"
	"lebron/apps/app/api/internal/types"
	"lebron/apps/user/rpc/user"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddReceiveAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// add user receiveAddress
func NewAddReceiveAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddReceiveAddressLogic {
	return &AddReceiveAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddReceiveAddressLogic) AddReceiveAddress(req *types.UserReceiveAddressAddReq) (resp *types.UserReceiveAddressAddRes, err error) {
	// todo: add your logic here and delete this line
	uid, err := l.ctx.Value("uid").(json.Number).Int64()

	var addRpcReq user.AddAddressRequest
	addRpcReq.Uid = uid

	copyErr := copier.Copy(&addRpcReq, req)
	if copyErr != nil {
		return nil, err
	}

	_, err = l.svcCtx.UserRPC.AddAddress(l.ctx, &addRpcReq)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

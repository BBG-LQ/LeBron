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

type GetReceiveAddressListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get user receiveAddress list
func NewGetReceiveAddressListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReceiveAddressListLogic {
	return &GetReceiveAddressListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReceiveAddressListLogic) GetReceiveAddressList(req *types.UserReceiveAddressListReq) (resp *types.UserReceiveAddressListRes, err error) {
	// todo: add your logic here and delete this line
	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	var getAddListreq user.ListAddressesRequest
	getAddListreq.Uid = uid

	respRpcs, err := l.svcCtx.UserRPC.ListAddresses(l.ctx, &getAddListreq)
	if err != nil {
		return nil, err
	}

	var addresslist []types.UserReceiveAddress
	for _, address := range respRpcs.Addresses {
		var addresstemp types.UserReceiveAddress
		_ = copier.Copy(&addresstemp, address)
		addresslist = append(addresslist, addresstemp)
	}

	return &types.UserReceiveAddressListRes{
		List: addresslist,
	}, nil
}

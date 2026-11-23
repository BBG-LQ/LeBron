package logic

import (
	"context"

	"lebron/apps/user/rpc/internal/svc"
	"lebron/apps/user/rpc/model"
	"lebron/apps/user/rpc/user"
	"lebron/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListAddressesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAddressesLogic {
	return &ListAddressesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取收货地址列表
func (l *ListAddressesLogic) ListAddresses(in *user.ListAddressesRequest) (*user.ListAddressesResponse, error) {
	// todo: add your logic here and delete this line

	receiveAddressesList, err := l.svcCtx.UserReceiveAddressModel.FindAllByUid(l.ctx, uint64(in.GetUid()))

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "Failed  get user's receive address list err : %v , in :%+v", err, in)
	}

	var resp []*user.UserReceiveAddress

	for _, receiveAddresses := range receiveAddressesList {
		var pbAddress user.UserReceiveAddress

		_ = copier.Copy(&pbAddress, receiveAddresses)

		resp = append(resp, &pbAddress)
	}

	return &user.ListAddressesResponse{
		Addresses: resp,
	}, nil
}

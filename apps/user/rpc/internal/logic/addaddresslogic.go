package logic

import (
	"context"

	"lebron/apps/user/rpc/internal/svc"
	"lebron/apps/user/rpc/model"
	"lebron/apps/user/rpc/user"
	"lebron/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAddressLogic {
	return &AddAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 添加收货地址
func (l *AddAddressLogic) AddAddress(in *user.AddAddressRequest) (*user.AddAddressResponse, error) {
	// todo: add your logic here and delete this line

	dbAddress := new(model.UserReceiveAddress)
	dbAddress.Uid = uint64(in.GetUid())
	dbAddress.Name = in.GetName()
	dbAddress.Phone = in.GetPhone()
	dbAddress.Province = in.GetProvince()
	dbAddress.City = in.GetCity()
	dbAddress.Region = in.GetRegion()
	dbAddress.DetailAddress = in.GetDetailAddress()
	dbAddress.PostCode = in.GetPostCode()
	dbAddress.IsDefaul = uint64(in.GetIsDefault())

	_, err := l.svcCtx.UserReceiveAddressModel.Insert(l.ctx, dbAddress)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "AddUserReceiveAddress Database Exception : %+v , err: %v", dbAddress, err)
	}

	return &user.AddAddressResponse{}, nil
}

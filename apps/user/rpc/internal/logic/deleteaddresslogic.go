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

type DeleteAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAddressLogic {
	return &DeleteAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除收货地址
func (l *DeleteAddressLogic) DeleteAddress(in *user.DeleteAddressRequest) (*user.DeleteAddressResponse, error) {
	// todo: add your logic here and delete this line

	_, err := l.svcCtx.UserReceiveAddressModel.FindOne(l.ctx, uint64(in.GetId()))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.Wrap(xerr.NewErrMsg("数据不存在"), "数据不存在")
		}
		return nil, err
	}

	dbAddress := new(model.UserReceiveAddress)
	dbAddress.Id = uint64(in.GetId())
	dbAddress.IsDelete = 1

	err = l.svcCtx.UserReceiveAddressModel.UpdateIsDelete(l.ctx, dbAddress)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "DelUserReceiveAddress Database Exception : %+v , err: %v", dbAddress, err)
	}

	return &user.DeleteAddressResponse{}, nil
}

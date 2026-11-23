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
	"google.golang.org/grpc/status"
)

type EditAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditAddressLogic {
	return &EditAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 编辑收货地址
func (l *EditAddressLogic) EditAddress(in *user.EditAddressRequest) (*user.EditAddressResponse, error) {
	// todo: add your logic here and delete this line

	_, err := l.svcCtx.UserReceiveAddressModel.FindOne(l.ctx, uint64(in.GetId()))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "数据不存在")
		}
		return nil, err
	}

	if in.IsDefault == 1 {

	}

	dbAddress := new(model.UserReceiveAddress)

	errCopy := copier.Copy(&dbAddress, in)
	if errCopy != nil {
		return nil, errCopy
	}

	err = l.svcCtx.UserReceiveAddressModel.Update(l.ctx, dbAddress)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "AddUserReceiveAddress Database Exception : %+v , err: %v", dbAddress, err)
	}

	return &user.EditAddressResponse{}, nil
}

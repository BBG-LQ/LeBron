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

type GetAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressLogic {
	return &GetAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据主键id,查询收获地址
func (l *GetAddressLogic) GetAddress(in *user.GetAddressRequest) (*user.GetAddressResponse, error) {
	// todo: add your logic here and delete this line

	ReceiveAddress, err := l.svcCtx.UserReceiveAddressModel.FindOne(l.ctx, uint64(in.Id))

	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.Wrap(xerr.NewErrMsg("收获地址数据不存在"), "收获地址数据不存在")

		}
		return nil, err
	}

	var resp user.UserReceiveAddress

	_ = copier.Copy(&resp, ReceiveAddress)
	resp.CreateTime = ReceiveAddress.CreateTime.Unix()
	resp.UpdateTime = ReceiveAddress.CreateTime.Unix()

	return &user.GetAddressResponse{
		Address: &resp,
	}, nil
}

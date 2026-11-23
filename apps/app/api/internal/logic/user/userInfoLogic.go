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

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get user info
func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	userInfo, err := l.svcCtx.UserRPC.GetUser(l.ctx, &user.GetUserRequest{
		// Id: uid,
		Id: uid,
	})

	if err != nil {
		return nil, err
	}

	var user types.UserInfo

	_ = copier.Copy(&user, userInfo.User)
	return &types.UserInfoResp{
		UserInfo: user,
	}, nil
}

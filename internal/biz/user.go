package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/viper"
	"gitlab.top.slotssprite.com/my/api-layout/party/auth"
	"gitlab.top.slotssprite.com/my/api-layout/party/grpcx/status"

	"gitlab.top.slotssprite.com/my/api-layout/internal/service"
	"gitlab.top.slotssprite.com/my/api-layout/party/httpx/ginx"
	"gitlab.top.slotssprite.com/my/api-layout/party/statusx"
)

// NewUserBiz ...
func NewUserBiz(userSvc *service.UserService) *UserBiz {
	return &UserBiz{userSvc: userSvc}
}

type LoginParams struct {
	Phone string `json:"phone"`
}

// UserBiz ...
type UserBiz struct {
	userSvc *service.UserService
}

// Login ...
func (s *UserBiz) Login(ctx *ginx.Context) {
	log.Context(ctx).Info("---------UserBiz.login---------")

	var login LoginParams
	if err := ctx.ShouldBindJSON(&login); err != nil {
		_err := status.FromErrorWithMsg(err)
		ctx.RenderResult(_err)
		return
	}

	resp, err := s.userSvc.GetAccountInfo(ctx, login.Phone)
	if err != nil {
		ctx.RenderResult(err)
		return
	}

	se := viper.GetString("jwt.secret")
	t := viper.GetInt32("jwt.tokenExpire")
	i := viper.GetString("jwt.issuer")
	su := viper.GetString("jwt.sub")
	token := auth.NewJwt(se, su, i, t, resp.Id)
	// NOTE: 实际情况按协议 render，一般不会直接 render pb struct
	data := map[string]interface{}{"user": resp, "token": token}
	ctx.Render(statusx.StatusOk, data)
}

func (s *UserBiz) GetUserInfo(ctx *ginx.Context) {
	log.Context(ctx).Info("---------UserBiz.GetUserInfo---------")

	resp, err := s.userSvc.GetPlayerInfo(ctx)
	if err != nil {
		ctx.RenderResult(err)
		return
	}

	// NOTE: 实际情况按协议 render，一般不会直接 render pb struct
	data := map[string]interface{}{"user": resp}
	ctx.Render(statusx.StatusOk, data)
}

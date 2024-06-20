package service

import (
	pbuser "gitlab.top.slotssprite.com/my/api-layout/api/helloworld/v1/user"
	"gitlab.top.slotssprite.com/my/api-layout/party/auth"
	"gitlab.top.slotssprite.com/my/api-layout/party/grpcx/status"
	"gitlab.top.slotssprite.com/my/api-layout/party/httpx"
	"gitlab.top.slotssprite.com/my/api-layout/party/httpx/ginx"
	"gitlab.top.slotssprite.com/my/api-layout/party/statusx"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserService ...
type UserService struct {
	userCli pbuser.UserServiceClient
}

// NewUserRpcClient ...
func NewUserRpcClient() pbuser.UserServiceClient {
	conn := GetRpcLayoutConn()
	return pbuser.NewUserServiceClient(conn)
}

// NewUserService ...
func NewUserService(userCli pbuser.UserServiceClient) *UserService {
	return &UserService{userCli: userCli}
}

// GetAccountInfo ...
func (u *UserService) GetAccountInfo(ctx *ginx.Context, phone string) (*pbuser.Account, *httpx.Result) {
	resp, err := u.userCli.GetAccount(ctx.Request.Context(), &pbuser.GetAccountReq{
		Phone: phone,
	})
	if err != nil {
		_err := status.FromErrorWithMsg(err)
		if status.IsNotFoundError(_err) {
			return nil, httpx.NewResult(statusx.StatusNotFound, err)
		}
		return nil, _err
	}
	return resp, nil
}

// GetPlayerInfo ...
func (u *UserService) GetPlayerInfo(ctx *ginx.Context) (*pbuser.Player, *httpx.Result) {
	userId := auth.GetJwtUser(ctx)
	if userId <= 0 {
		return nil, httpx.NewResult(statusx.StatusUnauthorizedUser, nil)
	}

	resp, err := u.userCli.GetPlayerById(ctx.Request.Context(), &emptypb.Empty{})
	if err != nil {
		_err := status.FromErrorWithMsg(err)
		if status.IsNotFoundError(_err) {
			return nil, httpx.NewResult(statusx.StatusNotFound, err)
		}
		return nil, _err
	}
	return resp, nil
}

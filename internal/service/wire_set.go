package service

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

// ProviderSet ...
var ProviderSet = wire.NewSet(NewUserRpcClient, NewUserService)

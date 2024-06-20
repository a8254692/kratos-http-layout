//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"gitlab.top.slotssprite.com/my/api-layout/internal/biz"
	"gitlab.top.slotssprite.com/my/api-layout/internal/router"
	"gitlab.top.slotssprite.com/my/api-layout/internal/server"
	"gitlab.top.slotssprite.com/my/api-layout/internal/service"
	"gitlab.top.slotssprite.com/my/api-layout/party/bootstrap"
)

// initApp init kratos application.
func initApp(*bootstrap.ServiceInfo, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(router.ProviderSet, server.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

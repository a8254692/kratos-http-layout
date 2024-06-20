package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/spf13/pflag"
	_ "go.uber.org/automaxprocs"

	"gitlab.top.slotssprite.com/my/api-layout/party/bootstrap"
	"gitlab.top.slotssprite.com/my/api-layout/party/env"
	"gitlab.top.slotssprite.com/my/api-layout/party/httpx"
)

func newApp(svc *bootstrap.ServiceInfo, logger log.Logger, hs *http.Server, ms *httpx.MonitorServer) *kratos.App {
	return kratos.New(
		kratos.ID(svc.GetInstanceId()),
		kratos.Name(svc.Name),
		kratos.Version(svc.Version),
		kratos.Metadata(svc.Metadata),
		kratos.Logger(logger),
		kratos.Server(hs, (*http.Server)(ms)),
	)
}

func main() {
	// 定义一个命令行Flag来指定配置文件路径
	path := pflag.StringP("config", "c", "", "Path to the config file")
	pflag.Parse()

	err := bootstrap.LoadConfig(path)
	if err != nil {
		panic("load config error: " + err.Error())
	}

	svc := bootstrap.NewServiceInfo(env.GetServiceName(), env.GetServiceVersion())
	if svc == nil {
		panic("bootstrap svc is nil")
	}

	_, err = bootstrap.NewTracerProvider(svc)
	if err != nil {
		panic("tracer is nil")
	}

	logger, logCleanup := bootstrap.NewLogger()
	if logger == nil {
		panic("logger is nil")
	}
	defer logCleanup()

	app, cleanup, err := initApp(svc, logger)
	if err != nil {
		panic("init app error: " + err.Error())
	}
	defer cleanup()

	if err = app.Run(); err != nil {
		panic(err)
	}
}

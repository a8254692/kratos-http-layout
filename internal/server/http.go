package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	KHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/spf13/viper"

	"gitlab.top.slotssprite.com/my/api-layout/internal/conf"
	"gitlab.top.slotssprite.com/my/api-layout/party/env"
	"gitlab.top.slotssprite.com/my/api-layout/party/httpx"
	"gitlab.top.slotssprite.com/my/api-layout/party/monitor"
	"gitlab.top.slotssprite.com/my/api-layout/party/util"
)

// NewHTTPServer ...
func NewHTTPServer(app *gin.Engine) (*KHttp.Server, error) {
	opts := make([]KHttp.ServerOption, 0)

	httpHost := viper.GetString(conf.PathHttpHost)
	httpPort := viper.GetInt(conf.PathHttpPort)
	httpTimeout := viper.GetInt(conf.PathHttpTimeout)

	if httpHost != "" && httpPort > 0 {
		opts = append(opts, KHttp.Address(fmt.Sprintf("%s:%d", httpHost, httpPort)))
	}
	if httpTimeout >= 0 {
		// NOTE: context 的超时时间
		opts = append(opts, KHttp.Timeout(time.Duration(httpTimeout)*time.Second))
	}

	httpSrv := KHttp.NewServer(opts...)
	httpSrv.HandlePrefix("/", app)
	return httpSrv, nil
}

// NewMonitorHTTPServer ...
func NewMonitorHTTPServer() (*httpx.MonitorServer, error) {
	opts := make([]KHttp.ServerOption, 0)
	httpHost := viper.GetString(env.PathMonitorHttpHost)
	httpPort := viper.GetInt(env.PathMonitorHttpPort)
	if httpPort < 0 {
		var err error
		if httpPort, err = util.GetFreePort(); err != nil {
			return nil, err
		}
	}
	if httpHost != "" && httpPort > 0 {
		opts = append(opts, KHttp.Address(fmt.Sprintf("%s:%d", httpHost, httpPort)))
	}
	httpSrv := KHttp.NewServer(opts...)
	httpSrv.HandlePrefix("/", monitor.DefaultServeMux)
	return (*httpx.MonitorServer)(httpSrv), nil
}

package bootstrap

import (
	"context"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"gitlab.top.slotssprite.com/my/api-layout/party/bootstrap/logx"
	"gitlab.top.slotssprite.com/my/api-layout/party/env"
	"gitlab.top.slotssprite.com/my/api-layout/party/runtimex"
	"gitlab.top.slotssprite.com/my/api-layout/party/util"
)

// NewLogger ...
func NewLogger() (logger klog.Logger, cleanup func()) {
	switch env.GetMode() {
	case env.ModeDevelop, env.ModeTest, env.ModeProduction:
		logger, cleanup = logx.NewLogrusLogger()
	case env.ModeLocal:
		logger, cleanup = logx.NewLogrusLogger()
	default:
		return nil, nil
	}

	time.Now()

	localIP, _ := util.GetLocalIP()
	logger = klog.With(logger,
		"@system", env.GetServiceName(),
		"@version", env.GetServiceVersion(),
		"@source", localIP,
		"@caller", runtimex.Caller(4),
		"@spanId", tracing.SpanID(),
		"@traceId", tracing.TraceID(),
		"@timestamp", _timestamp(),
		"@date", klog.Timestamp(time.DateTime),
	)
	return
}

// _timestamp ...
func _timestamp() klog.Valuer {
	return func(ctx context.Context) interface{} {
		return time.Now().Unix()
	}
}

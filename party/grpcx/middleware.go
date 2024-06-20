package grpcx

import (
	"context"

	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"gitlab.top.slotssprite.com/my/api-layout/party/grpcx/group"
	"gitlab.top.slotssprite.com/my/api-layout/party/grpcx/status"
	"gitlab.top.slotssprite.com/my/api-layout/party/statusx"
)

// ClientBreaker ...
func ClientBreaker() middleware.Middleware {
	gp := group.NewGroup(func() interface{} {
		// OPTIMIZE: NewBreaker Option ...
		return sre.NewBreaker()
	})
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			info, _ := transport.FromClientContext(ctx)
			breaker := gp.Get(info.Operation()).(circuitbreaker.CircuitBreaker)
			if err := breaker.Allow(); err != nil {
				breaker.MarkFailed()
				return nil, status.NewError(err, statusx.StatusTemporarilyUnavailable)
			}
			reply, err := handler(ctx, req)
			if err != nil && (errors.IsInternalServer(err) || errors.IsServiceUnavailable(err) || errors.IsGatewayTimeout(err)) {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
			return reply, err
		}
	}
}

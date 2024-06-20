package ginx

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	KHttp "github.com/go-kratos/kratos/v2/transport/http"

	"gitlab.top.slotssprite.com/my/api-layout/party/monitor"
	"gitlab.top.slotssprite.com/my/api-layout/party/statusx"
	"gitlab.top.slotssprite.com/my/api-layout/party/util"
)

// Middlewares return middlewares wrapper
func Middlewares(m ...middleware.Middleware) gin.HandlerFunc {
	chain := middleware.Chain(m...)
	return func(c *gin.Context) {
		next := func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			switch c.Writer.Status() {
			case http.StatusTooManyRequests:
				if ginCtx, ok := FromGinContext(ctx); ok {
					NewContext(ginCtx).Render(statusx.StatusTooManyRequests, nil)
				}
			}
			return
		}
		next = chain(next)
		ctx := NewGinContext(c.Request.Context(), c)
		if ginCtx, ok := FromGinContext(ctx); ok {
			KHttp.SetOperation(ctx, ginCtx.FullPath())
		}
		_, _ = next(ctx, c.Request)
	}
}

// Recovery ...
func Recovery() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			defer func() {
				if _err := recover(); _err != nil {
					buf := make([]byte, 64<<10)
					n := runtime.Stack(buf, false)
					errInfo := map[string]interface{}{
						"error": fmt.Sprintf("%v", _err),
						"req":   fmt.Sprintf("%+v", req),
						"stack": fmt.Sprintf("%s", buf[:n]),
					}
					log.Context(ctx).Error(util.MustMarshalToString(errInfo))
					if ginCtx, ok := FromGinContext(ctx); ok {
						NewContext(ginCtx).Render(statusx.StatusInternalServerError, nil, http.StatusInternalServerError)
						return
					}
				}
			}()
			return handler(ctx, req)
		}
	}
}

// Cors ...
func Cors() middleware.Middleware {
	corsHandlerFunc := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           3600 * time.Second,
	})
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if ginCtx, ok := FromGinContext(ctx); ok {
				if ginCtx.Request.Method == http.MethodOptions {
					corsHandlerFunc(ginCtx)
					return
				} else {
					ginCtx.Header("Access-Control-Allow-Origin", "*")
					ginCtx.Header("Access-Control-Allow-Credentials", "false")
				}
			}
			return handler(ctx, req)
		}
	}
}

// StartAt ...
func StartAt() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if ginCtx, ok := FromGinContext(ctx); ok {
				ginCtx.Set("startAt", time.Now().UnixMilli())
			}
			return handler(ctx, req)
		}
	}
}

// RateLimit ...
func RateLimit() middleware.Middleware {
	limiter := bbr.NewLimiter()
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			done, err := limiter.Allow()
			if err != nil {
				if ginCtx, ok := FromGinContext(ctx); ok {
					ginCtx.Status(http.StatusTooManyRequests)
					ginCtx.Abort()
					return handler(ctx, req)
				}
			}
			defer done(ratelimit.DoneInfo{Err: err})
			return handler(ctx, req)
		}
	}
}

var (
	prom     *monitor.Prom
	promOnce sync.Once
)

// Metrics ...
func Metrics(svcName string) middleware.Middleware {
	promOnce.Do(func() {
		labelNames := []string{"app", "method", "path", "code"}
		prom = monitor.NewProm("").
			RegisterCounter("http_request_handle_total", labelNames).
			RegisterHistogram("http_request_handle_seconds", labelNames)
	})
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			defer func() {
				if ginCtx, ok := FromGinContext(ctx); ok {
					labels := []string{svcName, ginCtx.Request.Method, ginCtx.FullPath(), strconv.Itoa(ginCtx.Writer.Status())}
					prom.CounterIncr(labels...)
					processTime := NewContext(ginCtx).processTime()
					prom.HistogramObserve(float64(processTime), labels...)
				}
			}()
			resp, err := handler(ctx, req)
			return resp, err
		}
	}
}

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"gitlab.top.slotssprite.com/my/api-layout/internal/biz"
	"gitlab.top.slotssprite.com/my/api-layout/party/env"
	"gitlab.top.slotssprite.com/my/api-layout/party/httpx/ginx"
)

// NewRouter ...
func NewRouter(user *biz.UserBiz) (app *gin.Engine) {
	app = gin.New()

	gin.SetMode(gin.DebugMode)

	// NOTE: 添加中间件注意执行顺序！
	//validate.Validator(),
	app.Use(ginx.Middlewares(ginx.StartAt(), tracing.Server(), ginx.Metrics(env.GetServiceName()), ginx.Recovery(), ginx.Cors(), ginx.RateLimit()))

	noCheck := app.Group("/api/v1")
	{
		noCheck.POST("/login", ginx.Handle(user.Login))
	}

	protected := app.Group("/api/v1", JwtAuth())
	{
		protected.POST("/user", ginx.Handle(user.GetUserInfo))
	}

	return
}

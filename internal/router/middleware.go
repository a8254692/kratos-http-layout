package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"gitlab.top.slotssprite.com/my/api-layout/party/util"
	"net/http"

	"gitlab.top.slotssprite.com/my/api-layout/internal/conf"
	"gitlab.top.slotssprite.com/my/api-layout/party/httpx/ginx"
	"gitlab.top.slotssprite.com/my/api-layout/party/statusx"
)

// JwtAuth ...
func JwtAuth() gin.HandlerFunc {
	secret := viper.GetString(conf.PathJwtSecret)
	secretKey := []byte(secret) // 替换为你的256位密钥
	return func(c *gin.Context) {
		// 从请求头中获取token字符串
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			ginx.NewContext(c).Render(statusx.StatusUnauthorized, "Unauthorized1", http.StatusUnauthorized)
			c.Abort()
			return
		}
		// 去掉Bearer前缀
		token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
			// 验证alg是否为HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil // 返回用于验证的密钥
		})
		if err != nil {
			ginx.NewContext(c).Render(statusx.StatusUnauthorized, "Unauthorized2", http.StatusUnauthorized)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 将claims信息附加到上下文，以便后续处理使用
			userId, ok := claims["user_id"].(float64)
			if !ok || userId <= 0 {
				ginx.NewContext(c).Render(statusx.StatusUnauthorized, "Unauthorized3", http.StatusUnauthorized)
				c.Abort()
				return
			}

			c.Set("jwt_token", claims)

			mc := metadata.AppendToClientContext(c.Request.Context(),
				"x-md-global-user-id", util.Int64ToString(int64(userId)),
			)

			c.Request = c.Request.WithContext(mc)
			c.Next()
		} else {
			ginx.NewContext(c).Render(statusx.StatusUnauthorized, "Unauthorized3", http.StatusUnauthorized)
			c.Abort()
			return
		}
	}
}

// UserAuth ...
//func UserAuth() gin.HandlerFunc {
//	userClient := service.NewUserRpcClient()
//	return func(c *gin.Context) {
//		var sessionId string
//		jwtToken, ok := c.Get("jwt_token")
//		if !ok {
//			ginx.NewContext(c).Render(statusx.StatusUnauthorized, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		jwtTokenInfo := jwtToken.(jwt.MapClaims)
//		if jwtTokenInfo != nil {
//			sessionId, _ = jwtTokenInfo["session_id"].(string)
//		}
//		if sessionId == "" {
//			ginx.NewContext(c).Render(statusx.StatusUnauthorized, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		user, err := userClient.GetPlayerById(ginx.NewContext(c), &emptypb.Empty{})
//		if err != nil {
//			_err := status.FromErrorWithMsg(err)
//			if status.IsNotFoundError(_err) {
//				ginx.NewContext(c).Render(statusx.StatusUnauthorized, "Unauthorized", http.StatusUnauthorized)
//				return
//			}
//			ginx.NewContext(c).RenderResult(_err)
//			return
//		}
//
//		// NOTE: metadata 用于记录日志
//		c.Set("metadata", map[string]interface{}{"session_id": sessionId, "user_id": user.Id})
//		c.Set("session_user", user)
//
//		mCtx := metadata.AppendToClientContext(c,
//			"x-md-auth-id", sessionId,
//			"x-md-auth-user-id", util.Int64ToString(user.Id),
//		)
//
//		c.Request = c.Request.WithContext(mCtx)
//		c.Next()
//	}
//}

package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"gitlab.top.slotssprite.com/my/api-layout/party/httpx/ginx"
)

// GetJwtUser ...
func GetJwtUser(ctx *ginx.Context) int64 {
	jwtToken, ok := ctx.Get("jwt_token")
	if !ok {
		return 0
	}

	jwtTokenInfo := jwtToken.(jwt.MapClaims)
	userId, ok := jwtTokenInfo["user_id"].(float64)
	if !ok {
		return 0
	}

	return int64(userId)
}

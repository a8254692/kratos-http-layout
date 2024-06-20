package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewJwt(secret, sub, name string, expTime int32, userId int64) string {
	// 定义密钥，确保这个密钥在验证JWT时也能访问到
	var signingKey = []byte(secret)

	// 定义一个 MapClaims 对象，存放自定义的数据声明
	claims := jwt.MapClaims{
		"sub":     sub,                                                       // 主题（Subject）
		"name":    name,                                                      // 名称
		"iat":     time.Now().Unix(),                                         // 签发时间
		"exp":     time.Now().Add(time.Hour * time.Duration(expTime)).Unix(), // 过期时间，这里设置为24小时后
		"user_id": userId,                                                    // 名称
	}

	// 创建一个 Token 对象，使用 HMAC 算法 HS256 进行签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥对Token进行签名
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		fmt.Println("Error generating JWT:", err)
		return ""
	}

	return tokenString
}

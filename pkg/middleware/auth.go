package middleware

import (
	"context"
	"douyin/pkg/database"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware 是一个中间件函数，用于校验 JWT 并检查黑名单
func AuthMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 从请求头中获取 JWT
		authHeader := string(ctx.GetHeader("Authorization")) // 将 []byte 转换为 string
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, utils.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		// 检查 Authorization 头格式是否正确
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, utils.H{"error": "Invalid Authorization header format"})
			ctx.Abort()
			return
		}

		// 提取 JWT
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 检查 JWT 是否在黑名单中
		val, err := database.RedisClient.Get(context.Background(), tokenString).Result()
		if err == nil && val == "blacklisted" {
			ctx.JSON(http.StatusUnauthorized, utils.H{"error": "token is blacklisted"})
			ctx.Abort()
			return
		}

		// 校验 JWT
		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil // 替换为你的 JWT 密钥
		})
		if err != nil || !parsedToken.Valid {
			ctx.JSON(http.StatusUnauthorized, utils.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		// 将解析后的用户 ID 存储到上下文中
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			ctx.Set("user_id", claims["user_id"])
		}

		// 继续执行后续的请求处理逻辑
		ctx.Next(c)
	}
}
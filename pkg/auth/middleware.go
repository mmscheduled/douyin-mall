package auth

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "douyin/pkg/config"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 检查白名单
        path := c.Request.URL.Path
        for _, whitePath := range config.AuthWhiteList {
            if path == whitePath {
                c.Next()
                return
            }
        }

        // 从 Header 中获取 Token
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
            return
        }

        // 解析 Token
        claims, err := ParseToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
            return
        }

        // 检查黑名单
        if _, ok := config.BlackList[claims.UserID]; ok {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "用户已被封禁"})
            return
        }

        // 将用户ID存入上下文
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
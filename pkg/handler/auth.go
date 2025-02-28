package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

func RefreshToken(c context.Context, ctx *app.RequestContext) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.H{"error": "invalid request"})
		return
	}

	// 解析刷新令牌
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_refresh_secret_key"), nil
	})
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, utils.H{"error": "invalid refresh token"})
		return
	}

	// 获取用户 ID
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to parse token claims"})
		return
	}
	userID := claims["user_id"].(float64)

	// 生成新的 JWT
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // 新的 JWT 过期时间：1 小时
	})
	newTokenString, err := newToken.SignedString([]byte("your_secret_key"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to generate new token"})
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"token": newTokenString,
	})
}
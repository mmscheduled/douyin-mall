package handler

import (
	"context"
	"douyin/pkg/database"
	"douyin/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func UserRegister(c context.Context, ctx *app.RequestContext) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.H{"error": "invalid request"})
		return
	}

	// 检查用户名是否已存在
	var user model.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err == nil {
		ctx.JSON(http.StatusConflict, utils.H{"error": "username already exists"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to hash password"})
		return
	}

	// 创建用户
	user = model.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to create user"})
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"user_id": user.ID,
		"token":   "dummy-token", // 实际项目中应生成 JWT
	})
}

func UserLogin(c context.Context, ctx *app.RequestContext) {
  var req struct {
      Username string `json:"username"`
      Password string `json:"password"`
  }
  if err := ctx.Bind(&req); err != nil {
      ctx.JSON(http.StatusBadRequest, utils.H{"error": "invalid request"})
      return
  }

  // 查找用户
  var user model.User
  if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
      ctx.JSON(http.StatusUnauthorized, utils.H{"error": "username or password is incorrect"})
      return
  }

  // 验证密码
  if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
      ctx.JSON(http.StatusUnauthorized, utils.H{"error": "username or password is incorrect"})
      return
  }

  // 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // JWT 过期时间：1 小时
	})
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to generate token"})
		return
	}

	// 生成刷新令牌
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 刷新令牌过期时间：7 天
	})
	refreshTokenString, err := refreshToken.SignedString([]byte("your_refresh_secret_key"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to generate refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"user_id":       user.ID,
		"token":         tokenString,
		"refresh_token": refreshTokenString,
	})
}

func GetUserInfo(c context.Context, ctx *app.RequestContext) {
  userID := ctx.Query("user_id")
  var user model.User
  if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
      ctx.JSON(http.StatusNotFound, utils.H{"error": "user not found"})
      return
  }

  ctx.JSON(http.StatusOK, utils.H{
      "username": user.Username,
      "email":    user.Email,
  })
}

func UserLogout(c context.Context, ctx *app.RequestContext) {
    token := string(ctx.GetHeader("Authorization")) // 从请求头中获取 JWT
    if token == "" {
        ctx.JSON(http.StatusBadRequest, utils.H{"error": "token is required"})
        return
    }

    // 将 JWT 加入 Redis 黑名单
    err := database.RedisClient.Set(context.Background(), token, "blacklisted", 0).Err()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to logout"})
        return
    }

    ctx.JSON(http.StatusOK, utils.H{
        "message": "logged out successfully",
    })
}

func DeleteUser(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Query("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, utils.H{"error": "user_id is required"})
		return
	}

	// 删除用户
	if err := database.DB.Where("id = ?", userID).Delete(&model.User{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"message": "user deleted successfully",
	})
}

func UpdateUser(c context.Context, ctx *app.RequestContext) {
	var req struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.H{"error": "invalid request"})
		return
	}

	// 更新用户信息
	if err := database.DB.Model(&model.User{}).Where("id = ?", req.UserID).Updates(map[string]interface{}{
		"username": req.Username,
		"email":    req.Email,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"message": "user updated successfully",
	})
}

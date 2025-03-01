package auth

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v4"
)

var (
    jwtSecret = []byte("your_secret_key") // JWT 密钥
)

// Claims JWT 数据结构
type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

// GenerateToken 生成 JWT
func GenerateToken(userID string) (string, error) {
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 有效期 24 小时
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// ParseToken 解析并验证 JWT
func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
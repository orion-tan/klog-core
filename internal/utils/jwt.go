package utils

import (
	"errors"
	"klog-backend/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type KLogClaims struct {
	UserID uint `json:"user_id"`
	Username string `json:"username"`
	Role string `json:"role"`
	Status string `json:"status"`
	jwt.RegisteredClaims
}

// GenerateToken 生成Token
// @userID 用户ID
// @username 用户名
// @role 角色
// @status 状态
// @return Token, 错误
func GenerateToken(userID uint, username, role, status string) (string, error) {
	claims := KLogClaims{
		UserID: userID,
		Username: username,
		Role: role,
		Status: status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Cfg.Jwt.ExpireHour) * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer: "klog-backend",
			Subject: username,
			Audience: jwt.ClaimStrings{"klog-backend"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Cfg.Jwt.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyToken 验证Token
// @token Token
// @return Claims, 错误
func VerifyToken(tokenString string) (*KLogClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &KLogClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*KLogClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token is invalid")
}
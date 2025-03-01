package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func LoadConfig() {

	viper.SetConfigName("config") // 配置文件名称

	viper.SetConfigType("yaml") // 配置文件类型

	viper.AddConfigPath(".") // 配置文件路径

	if err := viper.ReadInConfig(); err != nil {

		panic(fmt.Errorf("failed to read config file: %v", err))

	}

}

func GetDBConfig() string {

	return viper.GetString("database.dsn")

}

// 认证白名单（不需要认证的接口路径）
var AuthWhiteList = []string{
    "/api/user/login",
    "/api/user/register",
    "/health",
}

// 黑名单（用户ID列表）
var BlackList = map[string]bool{
    "blocked_user_1": true,
    "blocked_user_2": true,
}

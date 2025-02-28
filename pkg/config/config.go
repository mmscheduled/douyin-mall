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

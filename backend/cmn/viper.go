package cmn

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"runtime"
)

func init() {
	//指定配置文件路径
	sys := runtime.GOOS

	if sys == "windows" {
		viper.SetConfigFile("configure.json")
	} else if sys == "linux" {
		viper.SetConfigFile("/var/deploy/cst/backend/configure.json")
	}
	//查找并读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		zap.L().Error("配置文件出錯", zap.Error(err))
	}
}

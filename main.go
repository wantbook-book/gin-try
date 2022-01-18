package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"hannibal/gin-try/common"
	"os"
)

func main() {
	InitConfig()
	common.InitDB()
	//注册路由
	r := gin.Default()
	r = CollectRoute(r)
	//
	if port := viper.GetString("server.port"); port != "" {
		panic(r.Run(":" + port))
	}

}

func InitConfig() {
	//工作目录
	workDir, _ := os.Getwd()
	//配置文件名
	viper.SetConfigName("application")
	//配置文件类型
	viper.SetConfigType("yml")
	//配置文件目录
	viper.AddConfigPath(workDir + "/config")
	//读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read in config file")
	}
}

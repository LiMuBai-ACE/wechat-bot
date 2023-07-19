package main

import (
	"fmt"
	"io/ioutil"
	"wechat-bot/core"
	"wechat-bot/global"
	"wechat-bot/log"
	"wechat-bot/route"
	"wechat-bot/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func init() {

	// 读取配置文件
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %v", err))
	}

	// 解析YAML文件内容
	err = yaml.Unmarshal(data, &core.SystemConfig)
	if err != nil {
		panic(fmt.Sprintf("配置文件解析失败: %v", err))
	}

	// 手动初始化日志
	if gin.Mode() == gin.DebugMode {
		log.InitLogger(core.SystemConfig.LogConfig)
	}

	log.Debugf("配置文件解析完成: %v", core.SystemConfig)
}

// 程序启动入口
func main() {
	// 初始化Gin
	app := gin.Default()

	// 设置静态资源服务
	rootPath := utils.GetDirPath()
	for _, v := range core.SystemConfig.Serve.Static {
		app.Static(v, fmt.Sprintf("%s%s", rootPath, v))
	}

	// 定义全局异常处理
	app.NoRoute(core.NotFoundErrorHandler())
	// AppKey预检
	// app.Use(middleware.CheckAppKeyExistMiddleware(), middleware.CheckAppKeyIsLoggedInMiddleware())
	// 初始化路由
	route.InitRoute(app)

	// 初始化WechatBotMap
	global.InitWechatBotsMap()

	// 初始化Redis里登录的数据
	// global.InitBotWithStart()

	// 监听端口
	err := app.Run(core.SystemConfig.Serve.Host)
	if err != nil {
		panic(fmt.Sprintf("服务启动失败：%v", err))
	}
}

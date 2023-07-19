package route

import (
	"wechat-bot/controller"

	"github.com/gin-gonic/gin"
)

// initLoginRoute 初始化登录路由信息
func initLoginRoute(app *gin.Engine) {
	// 检查登录状态
	app.POST("/login", controller.LoginHandle)
}

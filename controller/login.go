package controller

import (
	"fmt"
	"io/fs"
	"os"
	"wechat-bot/core"
	"wechat-bot/global"
	log "wechat-bot/log"

	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
)

// LoginHandle 登录
func LoginHandle(ctx *gin.Context) {
	appKey := ctx.Request.Header.Get("AppKey")

	if appKey == "" {
		core.FailWithMessage("请求头 AppKey 为必传参数", ctx)
		return
	}

	//usePush := ctx.Query("usePush") // 是否使用免扫码登录
	//isPush := usePush == "1" || usePush == "true" || usePush == "yes"

	// 获取Bot对象
	bot := global.GetBot(appKey)
	if bot == nil {
		// 重新建bot对象 更新全局变量
		bot = global.InitWechatBotHandle()
		global.SetBot(appKey, bot)
	}

	// 本地创建数据缓存文件
	fileName := fmt.Sprintf("%s%s%s", "storage/", appKey, ".json")
	_, err := os.Stat(fileName)
	if err != nil {
		_, err := os.ReadDir("storage")
		if err != nil {
			os.Mkdir("storage", fs.ModePerm)
			os.Create(fileName)
		}
	}

	// 定义登录数据缓存
	storage := openwechat.NewFileHotReloadStorage(fileName)

	// 热登录
	var opts []openwechat.BotLoginOption
	opts = append(opts, openwechat.NewRetryLoginOption()) // 热登录失败使用扫码登录，适配第一次登录的时候无热登录数据
	//opts = append(opts, openwechat.NewSyncReloadDataLoginOption(10*time.Minute)) // 十分钟同步一次热登录数据

	// 登录
	if err := bot.HotLogin(storage, opts...); err != nil {
		log.Errorf("登录失败: %v", err)
		core.FailWithMessage("登录失败："+err.Error(), ctx)
		return
	}

	// 获取登录用户信息
	user, err := bot.GetCurrentUser()
	if err != nil {
		log.Errorf("获取登录用户信息失败: %v", err.Error())
		core.FailWithMessage("获取登录用户信息失败："+err.Error(), ctx)
		return
	}
	log.Infof("当前登录用户：%v", user.NickName)
	core.OkWithMessage("登录成功", ctx)
}

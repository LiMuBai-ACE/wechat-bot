package handler

import (
	"wechat-bot/plugins"

	"github.com/eatmoreapple/openwechat"
)

func HandleMessage(bot *openwechat.Bot) {
	// 定义一个处理器
	dispatcher := openwechat.NewMessageMatchDispatcher()
	// 设置为异步处理
	dispatcher.SetAsync(true)
	// 处理消息为已读
	// dispatcher.RegisterHandler(checkIsCanRead, setTheMessageAsRead)
	// 默认启用插件
	plugins.ChangePluginStatus(true)

	// 注册插件
	dispatcher.RegisterHandler(plugins.WeChatPluginInstance.CheckIsOpen, plugins.WeChatPluginInstance.Weather)

	// 注册文本消息处理函数
	dispatcher.OnText(textMessageHandle)

	// 注册图片消息处理器
	dispatcher.OnImage(imageMessageHandle)

	// 注册消息处理器
	bot.MessageHandler = dispatcher.AsMessageHandler()
}

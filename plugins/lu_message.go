package plugins

import (
	"regexp"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

func (weChatPlugin) LuMessage(ctx *openwechat.MessageContext) {

	if ctx.IsSendByGroup() {
		user, _ := ctx.Sender()
		group := &openwechat.Group{user}
		msg := ctx.Content
		if group.NickName == "中煤集团交流群" {
			// if msg == "八嘎" {
			// 	_, _ = ctx.ReplyText("八嘎雅璐！！ 撸几把出列！！")
			// } else
			if strings.Contains(msg, "骂") && strings.Contains(msg, "鲁") {
				re := regexp.MustCompile(`\d+`)
				numbers := re.FindAllString(msg, -1)
				if len(numbers) > 0 {
					// count, err := strconv.Atoi(numbers[0])
					// if err == nil {
					// 	for i := 0; i < count; i++ {
					// 		_, _ = ctx.ReplyText(fmt.Sprintf("鲁几把，是个阴阳蛋！%v", i+1))
					// 	}
					// }
					_, _ = ctx.ReplyText("八嘎，不许说我伟哥！😡😡😡😡😡")
				}
			} else if strings.Contains(msg, "伟") && strings.Contains(msg, "帅") {
				_, _ = ctx.ReplyText("湿了湿了，流水水")
				_, _ = ctx.ReplyText("伟哥，我爱你")
				_, _ = ctx.ReplyText("伟哥，我要给你生猴子么么😘😘😘😍😍😍😍")
			}
			// else if strings.Contains(msg, "本群最帅") {
			// 	_, _ = ctx.ReplyText("当然是我辣，人见人爱，花见花开的伟少啊 🤭😁🎈__中煤集团鲁总")

			// 	filepath := "D:/我的/微信图片_20230710180442.jpg" // 替换为实际的图片路径

			// 	file, err := os.Open(filepath)
			// 	if err != nil {
			// 		fmt.Println("无法打开文件:", err)
			// 		return
			// 	}
			// 	defer file.Close()

			// 	reader := io.Reader(file)
			// 	ctx.ReplyImage(reader)
			// }
		}
	}
	// if(senderUser.UserName)
	// _, _ = ctx.ReplyText("盖亚")
}

package handler

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"wechat-bot/log"

	"github.com/eatmoreapple/openwechat"
)

// ImageMessageData 图片消息结构体
type ImageMessageData struct {
	XMLName xml.Name `xml:"msg"`
	Img     struct {
		Text           string `xml:",chardata"`
		AesKey         string `xml:"aeskey,attr"`
		EnCryVer       string `xml:"encryver,attr"`
		CdnThumbAesKey string `xml:"cdnthumbaeskey,attr"`
		CdnThumbUrl    string `xml:"cdnthumburl,attr"`
		CdnThumbLength string `xml:"cdnthumblength,attr"`
		CdnThumbHeight string `xml:"cdnthumbheight,attr"`
		CdnThumbWidth  string `xml:"cdnthumbwidth,attr"`
		CdnMidHeight   string `xml:"cdnmidheight,attr"`
		CdnMidWidth    string `xml:"cdnmidwidth,attr"`
		CdnHdHeight    string `xml:"cdnhdheight,attr"`
		CdnHdWidth     string `xml:"cdnhdwidth,attr"`
		CdnMidImgUrl   string `xml:"cdnmidimgurl,attr"`
		Length         int64  `xml:"length,attr"`
		CdnBigImgUrl   string `xml:"cdnbigimgurl,attr"`
		HdLength       int64  `xml:"hdlength,attr"`
		Md5            string `xml:"md5,attr"`
	} `xml:"img"`
}

// 处理图片消息
func imageMessageHandle(ctx *openwechat.MessageContext) {
	sender, _ := ctx.Sender()
	senderUser := sender.NickName
	if ctx.IsSendByGroup() {
		// 取出消息在群里面的发送者
		senderInGroup, _ := ctx.SenderInGroup()
		senderUser = fmt.Sprintf("%v[%v]", senderInGroup.NickName, senderUser)
	}
	// 解析xml为结构体
	var data ImageMessageData
	if strings.HasPrefix(ctx.Content, "@") && !strings.Contains(ctx.Content, " ") {
		log.Debug("消息内容为图片资源ID，不解析为结构体")
	} else if err := xml.Unmarshal([]byte(ctx.Content), &data); err != nil {
		log.Errorf("消息解析失败: %v", err.Error())
		log.Debugf("发信人: %v ==> 原始内容: %v", senderUser, ctx.Content)
		return
	}

	log.Infof("[收到新图片消息] == 发信人：%v", senderUser)
	// 下载图片资源
	fileResp, err := ctx.GetFile()
	if err != nil {
		log.Errorf("图片下载失败: %v", err.Error())
		return
	}
	defer fileResp.Body.Close()
	imgFileByte, err := io.ReadAll(fileResp.Body)
	if err != nil {
		log.Errorf("图片读取错误: %v", err.Error())
		return
	} else {
		// 读取文件相关信息
		contentType := http.DetectContentType(imgFileByte)
		fileType := strings.Split(contentType, "/")[1]
		fileName := fmt.Sprintf("%v.%v", ctx.MsgId, fileType)
		if user, err := ctx.Bot().GetCurrentUser(); err == nil {
			uin := user.Uin
			fileName = fmt.Sprintf("%v/%v", uin, fileName)
		}

		// 存到本地文件夹内
		dirPath := "./images/"

		// 创建目录
		if _, err := os.Stat(dirPath); err != nil && os.IsNotExist(err) {
			if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
				log.Errorf("创建图片目录失败: %v", err.Error())
				return
			}
		}
		filePath := fmt.Sprintf("%v%v", dirPath, fileName)
		if err = ioutil.WriteFile(filePath, imgFileByte, os.ModePerm); err != nil {
			log.Errorf("图片保存失败: %v", err.Error())
			return
		} else {
			fileUrl := fmt.Sprintf("%v%v", dirPath, fileName)
			log.Infof("图片保存成功，图片链接: %v", fileUrl)
			ctx.Content = fileUrl
		}
	}
	ctx.Next()
}

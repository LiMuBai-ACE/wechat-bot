package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

type Province struct {
	Cities       []City `json:"citys"`
	ProvinceName string `json:"provinceName"`
}

type City struct {
	CityName string `json:"citysName"`
}

func loadProvinceData() ([]Province, error) {
	currentDir, _ := os.Getwd()
	// 拼接文件相对路径
	filePath := filepath.Join(currentDir, "provinces.json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var provinces []Province
	err = json.Unmarshal(data, &provinces)
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

// 字符串中是否有相同的文字
func checkDuplicateStrings(str1, str2 string) bool {
	matches := 0
	for _, char := range str1 {
		if strings.ContainsRune(str2, char) {
			matches++
			if matches >= 2 {
				return true
			}
		}
	}
	return false
}

func checkLocation(provinces []Province, location string) string {
	for _, province := range provinces {
		for _, city := range province.Cities {
			if checkDuplicateStrings(location, city.CityName) {
				return city.CityName
			}
		}
		if checkDuplicateStrings(location, province.ProvinceName) {
			return province.ProvinceName
		}
	}
	return ""
}

func getWeather(city string) {
	url := fmt.Sprintf("http://www.lpv4.cn:10000/api/weather/?city=%v", city) // 替换为目标接口的URL

	// 发送GET请求
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("发送GET请求失败：%s", err)
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("读取响应内容失败：%s", err)
		return
	}

	// 打印响应内容
	fmt.Println(string(body))
}

// 获取城市天气
func (weChatPlugin) Weather(ctx *openwechat.MessageContext) {

	if ctx.IsSendByGroup() {
		user, _ := ctx.Sender()
		group := &openwechat.Group{user}
		msg := ctx.Content
		if group.NickName == "中煤集团交流群" && strings.Contains(msg, "天气") {

			provinces, err := loadProvinceData()
			if err != nil {
				ctx.ReplyText("城市数据未能解析，请稍后重试，或联系管理员")
				return
			}

			// 去除查询词
			query := strings.Replace(msg, "查询", "", 1)

			// 去除天气词
			query = strings.Replace(query, "天气", "", 1)

			// 去除前后空格
			query = strings.TrimSpace(query)

			result := checkLocation(provinces, query)

			if result != "" {
				getWeather(result)
			} else {
				ctx.ReplyText(fmt.Sprintf("未找到与位置 '%v' 匹配的地点", msg))
			}
		}
	}

}

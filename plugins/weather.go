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
type WeatherInfo struct {
	Date      string       `json:"date"`
	Week      string       `json:"week"`
	Type      string       `json:"type"`
	Low       string       `json:"low"`
	High      string       `json:"high"`
	FengXiang string       `json:"fengxiang"`
	FengLi    string       `json:"fengli"`
	Night     NightWeather `json:"night"`
	Air       AirQuality   `json:"air"`
	Tip       string       `json:"tip"`
}

type NightWeather struct {
	Type      string `json:"type"`
	FengXiang string `json:"fengxiang"`
	FengLi    string `json:"fengli"`
}

type AirQuality struct {
	AQI      int    `json:"aqi"`
	AQILevel int    `json:"aqi_level"`
	AQIName  string `json:"aqi_name"`
	CO       string `json:"co"`
	NO2      string `json:"no2"`
	O3       string `json:"o3"`
	PM10     string `jsonpm10"`
	PM2_5    string `json:"pm2.5"`
	SO2      string `jsonso2"`
}
type Weather struct {
	Success bool        `json:"success"`
	City    string      `json:"city"`
	Info    WeatherInfo `json:"info"`
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

func getWeather(city string) (string, error) {
	url := fmt.Sprintf("http://www.lpv4.cn:10000/api/weather/?city=%v", city) // 替换为目标接口的URL

	// 发送GET请求
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("发送GET请求失败：%s", err)
		return "", nil
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应内容失败：%s", err)
	}
	var weather Weather
	err = json.Unmarshal([]byte(body), &weather)
	if err != nil {
		return "", fmt.Errorf("JSON转换失败：%s", err)
	}
	city, date, week, weatherType, lowTemp, highTemp, fengxiang, fengli, nightType, nightFengxiang, nightFengli, tip :=
		weather.City, weather.Info.Date, weather.Info.Week, weather.Info.Type, weather.Info.Low, weather.Info.High, weather.Info.FengXiang, weather.Info.FengLi, weather.Info.Night.Type, weather.Info.Night.FengXiang, weather.Info.Night.FengLi, weather.Info.Tip
	// ctx.ReplyText("城市数据未能解析，请稍后重试，或联系管理员")
	message := fmt.Sprintf("城市：%s\n日期：%s\n星期：%s\n天气：%s\n温度范围：%s ~ %s\n风向：%s\n风力：%s\n夜间气：%s\n夜间风向：%s\n夜间风力：%s\n提示：%s",
		city, date, week, weatherType, lowTemp, highTemp, fengxiang, fengli, nightType, nightFengxiang, nightFengli, tip)

	return message, nil

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
				msg, err := getWeather(result)
				if err != nil {
					ctx.ReplyText(err.Error())
					return
				}
				ctx.ReplyText(msg)
			} else {
				ctx.ReplyText(fmt.Sprintf("未找到与位置 '%v' 匹配的地点", msg))
			}
		}
	}

}

package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
	"wechat-bot/core"

	"github.com/eatmoreapple/openwechat"
)

type City struct {
	// 城市名称
	CityName string `json:"cityName"`
	// 城市编码
	AdCode string `json:"adCode"`
	// 城市区号
	CityCode string `json:"cityCode"`
}
type Forecast struct {
	Date           string `json:"date"`            // 日期
	Week           string `json:"week"`            // 星期几
	DayWeather     string `json:"dayweather"`      // 白天天气状况
	NightWeather   string `json:"nightweather"`    // 晚上天气状况
	DayTemp        string `json:"daytemp"`         // 白天温度
	NightTemp      string `json:"nighttemp"`       // 晚上温度
	DayWind        string `json:"daywind"`         // 白天风向
	NightWind      string `json:"nightwind"`       // 晚上风向
	DayPower       string `json:"daypower"`        // 白天风力
	NightPower     string `json:"nightpower"`      // 晚上风力
	DayTempFloat   string `json:"daytemp_float"`   // 白天温度（浮点数）
	NightTempFloat string `json:"nighttemp_float"` // 晚上温度（浮点数）
}

type CityForecast struct {
	City       string     `json:"city"`       // 城市名称
	Adcode     string     `json:"adcode"`     // 城市编码
	Province   string     `json:"province"`   // 省份名称
	ReportTime string     `json:"reporttime"` // 报告时间
	Casts      []Forecast `json:"casts"`      // 天气预报列表
}

type Live struct {
	Province         string `json:"province"`          // 省份名称
	City             string `json:"city"`              // 城市名称
	Adcode           string `json:"adcode"`            // 城市编码
	Weather          string `json:"weather"`           // 当前天气状况
	Temperature      string `json:"temperature"`       // 当前温度
	WindDirection    string `json:"winddirection"`     // 当前风向
	WindPower        string `json:"windpower"`         // 当前风力
	Humidity         string `json:"humidity"`          // 当前湿度
	ReportTime       string `json:"reporttime"`        // 报告时间
	TemperatureFloat string `json:"temperature_float"` // 当前温度（浮点数）
	HumidityFloat    string `json:"humidity_float"`    // 当前湿度（浮点数）
}

type Weather struct {
	Status    string         `json:"status"`              // 状态
	Count     string         `json:"count"`               // 数量
	Info      string         `json:"info"`                // 信息
	InfoCode  string         `json:"infocode"`            // 信息代码
	Forecasts []CityForecast `json:"forecasts,omitempty"` // 预报天气预报列表（可选）
	Lives     []Live         `json:"lives,omitempty"`     // 实时天气数据列表（可选）
}

// 气象类型
type Extensions string

const (
	// base:返回实况天气
	Base Extensions = "base"
	// all:返回预报天气
	All Extensions = "all"
)

func loadProvinceData() ([]City, error) {
	currentDir, _ := os.Getwd()
	// 拼接文件相对路径
	filePath := filepath.Join(currentDir, "provinces.json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var provinces []City
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
			if matches >= utf8.RuneCountInString(str1) {
				return true
			}
		}
	}
	return false
}

func checkLocation(Citys []City, location string) (City, error) {
	for _, city := range Citys {
		if checkDuplicateStrings(location, city.CityName) {
			return city, nil
		}
	}
	return City{}, fmt.Errorf("未找到城市信息！")
}

// 获取预计天气
func getWeather(city City, extensions Extensions) (Weather, error) {
	url := fmt.Sprintf("%v?key=%v&city=%v&extensions=%v", core.SystemConfig.Amap.WeatherApi, core.SystemConfig.Amap.Key, city.AdCode, extensions) // 替换为目标接口的URL
	var extensionsText = "实时天气"
	if extensions == All {
		extensionsText = "预报天气"
	}

	// 发送GET请求
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("发送GET请求失败：%s", err)
		return Weather{}, nil
	}
	defer response.Body.Close()
	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Weather{}, fmt.Errorf("读取%s响应内容失败：%s", extensionsText, err)
	}

	weather := Weather{
		Forecasts: make([]CityForecast, 0),
		Lives:     make([]Live, 0),
	}

	err = json.Unmarshal([]byte(body), &weather)

	if err != nil {
		return Weather{}, fmt.Errorf("JSON转换失败：%s", err)
	}

	return weather, nil
}

func getLiveWeatherInfo(live Live) string {
	info := "当前实时天气信息：\n"
	info += "省份名称：" + live.Province + "\n"
	info += "城市名称：" + live.City + "\n"
	// info += "城市编码：" + live.Adcode + "\n"
	info += "当前天气状况：" + live.Weather + "\n"
	info += "当前温度：" + live.Temperature + "℃\n"
	info += "当前风向：" + live.WindDirection + "\n"
	info += "当前风力：" + live.WindPower + "\n"
	info += "当前湿度：" + live.Humidity + "%\n"
	info += "报告时间：" + live.ReportTime
	return info
}
func getForecastWeatherInfo(forecast Forecast) string {

	weekMap := map[string]string{
		"1": "星期一",
		"2": "星期二",
		"3": "星期三",
		"4": "星期四",
		"5": "星期五",
		"6": "星期六",
		"7": "星期日",
	}

	info := ""
	info += "日期：" + forecast.Date + "，" + weekMap[forecast.Week] + "\n"
	info += "白天天气状况：" + forecast.DayWeather + "\n"
	info += "晚上天气状况：" + forecast.NightWeather + "\n"
	info += "白天温度：" + forecast.DayTemp + "\n"
	info += "晚上温度：" + forecast.NightTemp + "\n"
	info += "白天风向：" + forecast.DayWind + "\n"
	info += "晚上风向：" + forecast.NightWind + "\n"
	info += "白天风力：" + forecast.DayPower + "\n"
	info += "晚上风力：" + forecast.NightPower
	return info
}

// 合并天气数据 并进行汇报
// forecastsWeather 预报天气
// livesWeather 实时天气
func processData(forecastsWeather, livesWeather Weather) (string, error) {
	var liveText = ""
	for _, live := range livesWeather.Lives {
		liveText += getLiveWeatherInfo(live)
	}

	// var forecastText = ""
	// for _, cityForecast := range forecastsWeather.Forecasts {
	// 	forecastText += cityForecast.City + "今天及未来三天的天气预报如下:"
	// 	for _, forecast := range cityForecast.Casts {
	// 		forecastText += "\n" + getForecastWeatherInfo(forecast)
	// 	}
	// }

	// if liveText != "" && forecastText != "" {
	// 	return liveText + "\n" + forecastText, nil
	// }
	if liveText != "" {
		return liveText, nil
	}
	return "", fmt.Errorf("合并天气数据失败")
}

// 获取城市天气
func (weChatPlugin) Weather(ctx *openwechat.MessageContext) {

	if ctx.IsSendByGroup() {
		user, _ := ctx.Sender()
		group := &openwechat.Group{user}
		msg := ctx.Content
		if group.NickName == "中煤集团交流群" && strings.Contains(msg, "天气") && strings.Contains(msg, "查询") {

			provinces, err := loadProvinceData()
			if err != nil {
				ctx.ReplyText(fmt.Sprintf("城市数据未能解析，请稍后重试，或联系管理员:%v", err.Error()))
				return
			}

			// 去除查询词
			query := strings.Replace(msg, "查询", "", 1)

			// 去除天气词
			query = strings.Replace(query, "天气", "", 1)

			// 去除前后空格
			query = strings.TrimSpace(query)

			city, err := checkLocation(provinces, query)

			if err == nil {
				forecastsWeather, err1 := getWeather(city, All)
				livesWeather, err2 := getWeather(city, Base)
				if err1 != nil {
					ctx.ReplyText(err.Error())
					return
				}
				if err2 != nil {
					ctx.ReplyText(err.Error())
					return
				}
				msg, err = processData(forecastsWeather, livesWeather)
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

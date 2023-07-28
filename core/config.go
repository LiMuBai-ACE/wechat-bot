package core

import "wechat-bot/log"

// SystemConfig 系统配置
var SystemConfig systemConfig

// 系统配置
type systemConfig struct {
	Serve     ServeConfig   `yaml:"serve"`
	LogConfig log.LogConfig `yaml:"log"`
	Amap      AmapConfig    `yaml:"amap"`
}

// 服务配置
type ServeConfig struct {
	Host string `yaml:"host"`
	// 静态文件地址
	Static []string `yaml:"static"`
}

// 环境
type mode string

const (
	Dev  mode = "development"
	Prod mode = "production"
)

// 高德地图配置
type AmapConfig struct {
	Key string `yaml:"key"`
	// 静态文件地址
	WeatherApi string `yaml:"weatherApi"`
}

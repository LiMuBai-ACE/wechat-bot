package core

import "wechat-bot/log"

// SystemConfig 系统配置
var SystemConfig systemConfig

// 系统配置
type systemConfig struct {
	Serve     ServeConfig   `yaml:"serve"`
	LogConfig log.LogConfig `yaml:"log"`
}

type ServeConfig struct {
	Host   string   `yaml:"host"`
	Static []string `yaml:"static"`
}

type mode string

const (
	Dev  mode = "development"
	Prod mode = "production"
)

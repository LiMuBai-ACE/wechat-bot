package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var config LogConfig

// InitLogger 初始化日志工具
func InitLogger(c LogConfig) {
	config = c
	var cores []zapcore.Core
	// 生成输出到控制台的Core
	consoleCore := initConsoleCore()
	cores = append(cores, consoleCore)

	// 输出到文件的Core
	if config.FileEnable {
		fileCore := initFileCore()
		cores = append(cores, fileCore)
	}

	// 增加 caller 信息
	// AddCallerSkip 输出的文件名和行号是调用封装函数的位置，而不是调用日志函数的位置
	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(logger)
}

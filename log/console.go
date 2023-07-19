package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 初始化打印到控制台的ZapCore
func initConsoleCore() zapcore.Core {
	// 配置 sugaredLogger
	writer := zapcore.AddSync(os.Stdout)

	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("[%v]", t.Format("2006-01-02 15:04:05.000")))
	}

	// 格式相关的配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 修改时间戳的格式
	encoderConfig.EncodeTime = customTimeEncoder
	// 日志级别使用大写带颜色显示
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 设置日志等级，如果是release模式，控制台只打印Error级别以上的日志
	logLevel := zapcore.DebugLevel
	if config.Mode == Prod {
		logLevel = zapcore.ErrorLevel
	}
	// 设置日志级别
	core := zapcore.NewCore(encoder, writer, logLevel)
	return core
}

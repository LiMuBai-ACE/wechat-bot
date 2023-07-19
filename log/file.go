package log

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 初始化LokiCore，使日志可以输出到文件
func initFileCore() zapcore.Core {
	hostname, _ := os.Hostname()
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/runtime-%s.log", hostname), // 日志文件的位置
		MaxSize:    10,                                           // 最大10M
		MaxBackups: 5,                                            // 保留旧文件的最大个数
		MaxAge:     30,                                           // 保留旧文件的最大天数
		Compress:   false,                                        // 是否压缩/归档旧文件
	}
	// 配置 sugaredLogger
	writer := zapcore.AddSync(lumberJackLogger)

	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("[%v]", t.Format("2006-01-02 15:04:05.000")))
	}

	// 格式相关的配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 修改时间戳的格式
	encoderConfig.EncodeTime = customTimeEncoder
	// 日志级别使用大写显示
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
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

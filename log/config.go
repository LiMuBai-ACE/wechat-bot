package log

type mode string

const (
	Dev  mode = "development"
	Prod mode = "production"
)

// LogConfig 日志配置
type LogConfig struct {
	Mode       mode `yaml:"mode"`       // 模式
	FileEnable bool `yaml:"fileEnable"` // 模式
}
